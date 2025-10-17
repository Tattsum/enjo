'use client'

import React, { useState } from 'react'
import { useMutation } from '@apollo/client'
import {
  GENERATE_IMAGE,
  GenerateImageData,
  GenerateImageVariables,
  ImageStyle,
  AspectRatio,
} from '@/lib/graphql/queries'

interface ImageGeneratorProps {
  inflammatoryText: string
  onImageGenerated?: (imageUrl: string) => void
}

const ImageGenerator: React.FC<ImageGeneratorProps> = ({
  inflammatoryText,
  onImageGenerated,
}) => {
  const [style, setStyle] = useState<ImageStyle>(ImageStyle.MEME)
  const [aspectRatio] = useState<AspectRatio>(AspectRatio.SQUARE)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  const [generateImage, { loading }] = useMutation<GenerateImageData, GenerateImageVariables>(
    GENERATE_IMAGE,
    {
      onCompleted: (data) => {
        setErrorMessage(null)
        if (onImageGenerated) {
          onImageGenerated(data.generateImage.imageUrl)
        }
      },
      onError: (error) => {
        setErrorMessage(`エラー: ${error.message}`)
      },
    }
  )

  const handleGenerateImage = () => {
    if (!inflammatoryText) return

    generateImage({
      variables: {
        input: {
          text: inflammatoryText,
          style,
          aspectRatio,
        },
      },
    })
  }

  return (
    <div className="space-y-4">
      {/* Style Selector */}
      <div className="flex items-center gap-4">
        <label htmlFor="image-style" className="text-sm font-medium text-gray-700">
          スタイル:
        </label>
        <select
          id="image-style"
          value={style}
          onChange={(e) => setStyle(e.target.value as ImageStyle)}
          className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-fire-500 text-gray-900 bg-white"
        >
          <option value={ImageStyle.MEME}>ミーム風</option>
          <option value={ImageStyle.REALISTIC}>リアル調</option>
          <option value={ImageStyle.ILLUSTRATION}>イラスト調</option>
          <option value={ImageStyle.DRAMATIC}>ドラマチック</option>
        </select>
      </div>

      {/* Generate Button */}
      <button
        onClick={handleGenerateImage}
        disabled={!inflammatoryText || loading}
        className="w-full px-6 py-3 bg-gradient-to-r from-fire-500 to-fire-600 text-white rounded-lg hover:from-fire-600 hover:to-fire-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        aria-label="画像を生成"
      >
        {loading ? (
          <>
            <span className="animate-spin">⏳</span>
            <span>生成中...</span>
          </>
        ) : (
          <>
            <span>🎨</span>
            <span>画像を生成</span>
          </>
        )}
      </button>

      {/* Error Message */}
      {errorMessage && (
        <div className="bg-red-50 border-l-4 border-red-400 p-4 rounded">
          <p className="text-sm text-red-700">{errorMessage}</p>
        </div>
      )}
    </div>
  )
}

export default ImageGenerator
