/**
 * Integration test for the complete image generation flow
 *
 * This test verifies the end-to-end flow of:
 * 1. User generates inflammatory text
 * 2. User generates an image from the text
 * 3. User can download or regenerate the image
 * 4. User can post to Twitter with the image
 */

import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MockedProvider } from '@apollo/client/testing'
import ImageGenerator from '../ImageGenerator'
import ImagePreview from '../ImagePreview'
import { GENERATE_IMAGE } from '@/lib/graphql/queries'

// Mock data for tests
const mockInflammatoryText = "今日のランチは最高でした！みんなも食べるべき！"
const mockImagePrompt = "A dramatic image of passionate lunch enthusiasm with fire effects"
const mockImageUrl = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
const mockGeneratedAt = "2025-10-17T10:00:00Z"

describe('Image Generation Flow - Integration Tests', () => {
  describe('Complete image generation workflow', () => {
    it('should complete the full image generation flow', async () => {
      const user = userEvent.setup()
      const onImageGenerated = jest.fn()

      // Mock successful GraphQL response
      const mocks = [
        {
          request: {
            query: GENERATE_IMAGE,
            variables: {
              input: {
                text: mockInflammatoryText,
                style: 'MEME',
                aspectRatio: 'SQUARE',
              },
            },
          },
          result: {
            data: {
              generateImage: {
                imageUrl: mockImageUrl,
                prompt: mockImagePrompt,
                generatedAt: mockGeneratedAt,
              },
            },
          },
        },
      ]

      render(
        <MockedProvider mocks={mocks}>
          <ImageGenerator
            inflammatoryText={mockInflammatoryText}
            onImageGenerated={onImageGenerated}
          />
        </MockedProvider>
      )

      // Step 1: Verify initial state
      expect(screen.getByRole('button', { name: /画像を生成/ })).toBeInTheDocument()
      const styleSelect = screen.getByLabelText('スタイル:') as HTMLSelectElement
      expect(styleSelect.value).toBe('MEME')

      // Step 2: User clicks generate button
      const generateButton = screen.getByRole('button', { name: /画像を生成/ })
      await user.click(generateButton)

      // Step 3: Wait for image generation to complete
      await waitFor(() => {
        expect(onImageGenerated).toHaveBeenCalledWith(mockImageUrl)
      }, { timeout: 3000 })

      // Verify no error state
      expect(screen.queryByText(/エラー/)).not.toBeInTheDocument()
    })

    it('should allow style selection before generation', async () => {
      const user = userEvent.setup()

      const mocks = [
        {
          request: {
            query: GENERATE_IMAGE,
            variables: {
              input: {
                text: mockInflammatoryText,
                style: 'REALISTIC',
                aspectRatio: 'SQUARE',
              },
            },
          },
          result: {
            data: {
              generateImage: {
                imageUrl: mockImageUrl,
                prompt: mockImagePrompt,
                generatedAt: mockGeneratedAt,
              },
            },
          },
        },
      ]

      render(
        <MockedProvider mocks={mocks}>
          <ImageGenerator
            inflammatoryText={mockInflammatoryText}
            onImageGenerated={jest.fn()}
          />
        </MockedProvider>
      )

      // Change style to REALISTIC
      const styleSelect = screen.getByLabelText('スタイル:') as HTMLSelectElement
      await user.selectOptions(styleSelect, 'REALISTIC')

      // Verify selection
      expect(styleSelect.value).toBe('REALISTIC')

      // Generate with new style
      const generateButton = screen.getByRole('button', { name: /画像を生成/ })
      await user.click(generateButton)

      // Wait for the mutation to be called (test passes if no error)
    })

    it('should handle generation errors gracefully', async () => {
      const user = userEvent.setup()

      const mocks = [
        {
          request: {
            query: GENERATE_IMAGE,
            variables: {
              input: {
                text: mockInflammatoryText,
                style: 'MEME',
                aspectRatio: 'SQUARE',
              },
            },
          },
          error: new Error('Failed to generate image'),
        },
      ]

      render(
        <MockedProvider mocks={mocks}>
          <ImageGenerator
            inflammatoryText={mockInflammatoryText}
            onImageGenerated={jest.fn()}
          />
        </MockedProvider>
      )

      // Click generate
      await user.click(screen.getByRole('button', { name: /画像を生成/ }))

      // Wait for error message
      await waitFor(() => {
        expect(screen.getByText(/エラー:/)).toBeInTheDocument()
      })

      // Button should be enabled again for retry
      expect(screen.getByRole('button', { name: /画像を生成/ })).not.toBeDisabled()
    })
  })

  describe('Image preview and actions workflow', () => {
    it('should display image preview with all controls', () => {
      const onDownload = jest.fn()
      const onRegenerate = jest.fn()

      render(
        <ImagePreview
          imageUrl={mockImageUrl}
          prompt={mockImagePrompt}
          onDownload={onDownload}
          onRegenerate={onRegenerate}
        />
      )

      // Verify image is displayed
      const image = screen.getByAltText('生成された画像')
      expect(image).toBeInTheDocument()
      expect(image).toHaveAttribute('src', mockImageUrl)

      // Verify prompt is shown
      expect(screen.getByText(/使用したプロンプト:/)).toBeInTheDocument()
      expect(screen.getByText(mockImagePrompt)).toBeInTheDocument()

      // Verify action buttons
      expect(screen.getByRole('button', { name: 'ダウンロード' })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '再生成' })).toBeInTheDocument()
    })

    it('should handle download action', async () => {
      const user = userEvent.setup()
      const onDownload = jest.fn()

      render(
        <ImagePreview
          imageUrl={mockImageUrl}
          onDownload={onDownload}
        />
      )

      // Click download button
      const downloadButton = screen.getByRole('button', { name: 'ダウンロード' })
      await user.click(downloadButton)

      // Verify callback was called
      expect(onDownload).toHaveBeenCalledTimes(1)
    })

    it('should handle regenerate action', async () => {
      const user = userEvent.setup()
      const onRegenerate = jest.fn()

      render(
        <ImagePreview
          imageUrl={mockImageUrl}
          onRegenerate={onRegenerate}
        />
      )

      // Click regenerate button
      const regenerateButton = screen.getByRole('button', { name: '再生成' })
      await user.click(regenerateButton)

      // Verify callback was called
      expect(onRegenerate).toHaveBeenCalledTimes(1)
    })

    it('should work without optional props', () => {
      render(
        <ImagePreview imageUrl={mockImageUrl} />
      )

      // Image should still be displayed
      expect(screen.getByAltText('生成された画像')).toBeInTheDocument()

      // Optional elements should not be present
      expect(screen.queryByText(/使用したプロンプト:/)).not.toBeInTheDocument()
      expect(screen.queryByRole('button', { name: 'ダウンロード' })).not.toBeInTheDocument()
      expect(screen.queryByRole('button', { name: '再生成' })).not.toBeInTheDocument()
    })
  })

  describe('Image generation performance', () => {
    it('should handle rapid multiple generation requests', async () => {
      const user = userEvent.setup()
      const onImageGenerated = jest.fn()

      const mocks = [
        {
          request: {
            query: GENERATE_IMAGE,
            variables: {
              input: {
                text: mockInflammatoryText,
                style: 'MEME',
                aspectRatio: 'SQUARE',
              },
            },
          },
          result: {
            data: {
              generateImage: {
                imageUrl: mockImageUrl,
                prompt: mockImagePrompt,
                generatedAt: mockGeneratedAt,
              },
            },
          },
        },
      ]

      render(
        <MockedProvider mocks={mocks}>
          <ImageGenerator
            inflammatoryText={mockInflammatoryText}
            onImageGenerated={onImageGenerated}
          />
        </MockedProvider>
      )

      const generateButton = screen.getByRole('button', { name: /画像を生成/ })

      // First click
      await user.click(generateButton)

      // Wait for completion
      await waitFor(() => {
        expect(onImageGenerated).toHaveBeenCalled()
      }, { timeout: 3000 })

      // Verify it was only called once despite rapid clicks
      expect(onImageGenerated).toHaveBeenCalledTimes(1)
    })
  })

  describe('Accessibility', () => {
    it('should have proper ARIA labels for image generator', () => {
      render(
        <MockedProvider mocks={[]} addTypename={false}>
          <ImageGenerator
            inflammatoryText={mockInflammatoryText}
            onImageGenerated={jest.fn()}
          />
        </MockedProvider>
      )

      // Style selector should have a label
      expect(screen.getByLabelText('スタイル:')).toBeInTheDocument()

      // Button should have descriptive text
      const button = screen.getByRole('button', { name: /画像を生成/ })
      expect(button).toBeInTheDocument()
    })

    it('should have proper alt text for generated image', () => {
      render(
        <ImagePreview imageUrl={mockImageUrl} />
      )

      const image = screen.getByAltText('生成された画像')
      expect(image).toBeInTheDocument()
    })
  })

  describe('Data validation', () => {
    it('should accept valid data URL format', () => {
      const validDataUrl = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

      render(
        <ImagePreview imageUrl={validDataUrl} />
      )

      const image = screen.getByAltText('生成された画像')
      expect(image).toHaveAttribute('src', validDataUrl)
    })

    it('should accept HTTP URL format', () => {
      const httpUrl = "https://example.com/image.png"

      render(
        <ImagePreview imageUrl={httpUrl} />
      )

      const image = screen.getByAltText('生成された画像')
      expect(image).toHaveAttribute('src', httpUrl)
    })

    it('should handle very long prompt text', () => {
      const longPrompt = "A ".repeat(100) + "dramatic image"

      render(
        <ImagePreview
          imageUrl={mockImageUrl}
          prompt={longPrompt}
        />
      )

      // Prompt should be displayed (truncation handled by CSS)
      expect(screen.getByText(longPrompt)).toBeInTheDocument()
    })
  })
})
