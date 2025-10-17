import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import ImagePreview from '../ImagePreview'

describe('ImagePreview', () => {
  const mockImageUrl = 'data:image/png;base64,mockImageData'
  const mockPrompt = 'A meme-style image with fire elements'
  const mockOnDownload = jest.fn()
  const mockOnRegenerate = jest.fn()

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('renders the image', () => {
    render(<ImagePreview imageUrl={mockImageUrl} />)

    const img = screen.getByRole('img', { name: /生成された画像/i })
    expect(img).toBeInTheDocument()
    expect(img).toHaveAttribute('src', mockImageUrl)
  })

  it('renders the prompt when provided', () => {
    render(<ImagePreview imageUrl={mockImageUrl} prompt={mockPrompt} />)

    expect(screen.getByText(/使用したプロンプト/i)).toBeInTheDocument()
    expect(screen.getByText(mockPrompt)).toBeInTheDocument()
  })

  it('does not render prompt section when prompt is not provided', () => {
    render(<ImagePreview imageUrl={mockImageUrl} />)

    expect(screen.queryByText(/使用したプロンプト/i)).not.toBeInTheDocument()
  })

  it('renders download button when onDownload is provided', () => {
    render(<ImagePreview imageUrl={mockImageUrl} onDownload={mockOnDownload} />)

    const downloadButton = screen.getByRole('button', { name: /ダウンロード/i })
    expect(downloadButton).toBeInTheDocument()
  })

  it('calls onDownload when download button is clicked', () => {
    render(<ImagePreview imageUrl={mockImageUrl} onDownload={mockOnDownload} />)

    const downloadButton = screen.getByRole('button', { name: /ダウンロード/i })
    fireEvent.click(downloadButton)

    expect(mockOnDownload).toHaveBeenCalledTimes(1)
  })

  it('renders regenerate button when onRegenerate is provided', () => {
    render(<ImagePreview imageUrl={mockImageUrl} onRegenerate={mockOnRegenerate} />)

    const regenerateButton = screen.getByRole('button', { name: /再生成/i })
    expect(regenerateButton).toBeInTheDocument()
  })

  it('calls onRegenerate when regenerate button is clicked', () => {
    render(<ImagePreview imageUrl={mockImageUrl} onRegenerate={mockOnRegenerate} />)

    const regenerateButton = screen.getByRole('button', { name: /再生成/i })
    fireEvent.click(regenerateButton)

    expect(mockOnRegenerate).toHaveBeenCalledTimes(1)
  })

  it('renders both buttons when both callbacks are provided', () => {
    render(
      <ImagePreview
        imageUrl={mockImageUrl}
        onDownload={mockOnDownload}
        onRegenerate={mockOnRegenerate}
      />
    )

    expect(screen.getByRole('button', { name: /ダウンロード/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /再生成/i })).toBeInTheDocument()
  })

  it('does not render buttons when callbacks are not provided', () => {
    render(<ImagePreview imageUrl={mockImageUrl} />)

    expect(screen.queryByRole('button', { name: /ダウンロード/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /再生成/i })).not.toBeInTheDocument()
  })
})
