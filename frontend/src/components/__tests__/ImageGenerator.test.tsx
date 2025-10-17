import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import ImageGenerator from '../ImageGenerator'
import { GENERATE_IMAGE } from '@/lib/graphql/queries'

describe('ImageGenerator', () => {
  const mockOnImageGenerated = jest.fn()
  const inflammatoryText = 'これは炎上しやすいテキストです'

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('renders the generate image button', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <ImageGenerator
          inflammatoryText={inflammatoryText}
          onImageGenerated={mockOnImageGenerated}
        />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /画像を生成/i })
    expect(button).toBeInTheDocument()
  })

  it('renders style selector', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <ImageGenerator
          inflammatoryText={inflammatoryText}
          onImageGenerated={mockOnImageGenerated}
        />
      </MockedProvider>
    )

    expect(screen.getByLabelText(/スタイル/i)).toBeInTheDocument()
  })

  it('generates image when button is clicked', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_IMAGE,
          variables: {
            input: {
              text: inflammatoryText,
              style: 'MEME',
              aspectRatio: 'SQUARE',
            },
          },
        },
        result: {
          data: {
            generateImage: {
              imageUrl: 'data:image/png;base64,mockImageData',
              prompt: 'A meme-style image',
              generatedAt: '2025-10-17T12:00:00Z',
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <ImageGenerator
          inflammatoryText={inflammatoryText}
          onImageGenerated={mockOnImageGenerated}
        />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /画像を生成/i })
    fireEvent.click(button)

    // Loading state
    await waitFor(() => {
      expect(screen.getByText(/生成中/i)).toBeInTheDocument()
    })

    // Image generated
    await waitFor(() => {
      expect(mockOnImageGenerated).toHaveBeenCalledWith(
        'data:image/png;base64,mockImageData'
      )
    })
  })

  it('displays error message when image generation fails', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_IMAGE,
          variables: {
            input: {
              text: inflammatoryText,
              style: 'MEME',
              aspectRatio: 'SQUARE',
            },
          },
        },
        error: new Error('Image generation failed'),
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <ImageGenerator
          inflammatoryText={inflammatoryText}
          onImageGenerated={mockOnImageGenerated}
        />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /画像を生成/i })
    fireEvent.click(button)

    await waitFor(() => {
      expect(screen.getByText(/エラー/i)).toBeInTheDocument()
    })
  })

  it('disables button when text is empty', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <ImageGenerator inflammatoryText="" onImageGenerated={mockOnImageGenerated} />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /画像を生成/i })
    expect(button).toBeDisabled()
  })

  it('allows style selection', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <ImageGenerator
          inflammatoryText={inflammatoryText}
          onImageGenerated={mockOnImageGenerated}
        />
      </MockedProvider>
    )

    const styleSelect = screen.getByLabelText(/スタイル/i) as HTMLSelectElement
    fireEvent.change(styleSelect, { target: { value: 'REALISTIC' } })

    expect(styleSelect.value).toBe('REALISTIC')
  })
})
