import React from 'react'

interface ImagePreviewProps {
  imageUrl: string
  prompt?: string
  onDownload?: () => void
  onRegenerate?: () => void
}

const ImagePreview: React.FC<ImagePreviewProps> = ({
  imageUrl,
  prompt,
  onDownload,
  onRegenerate,
}) => {
  return (
    <div className="space-y-4">
      {/* Image Display */}
      <div className="relative rounded-lg overflow-hidden border-2 border-fire-300 shadow-lg">
        <img
          src={imageUrl}
          alt="生成された画像"
          className="w-full h-auto object-contain bg-gray-100"
        />
      </div>

      {/* Prompt Display */}
      {prompt && (
        <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
          <h4 className="text-sm font-semibold text-gray-700 mb-2">使用したプロンプト:</h4>
          <p className="text-sm text-gray-600">{prompt}</p>
        </div>
      )}

      {/* Action Buttons */}
      {(onDownload || onRegenerate) && (
        <div className="flex gap-3">
          {onRegenerate && (
            <button
              onClick={onRegenerate}
              className="flex-1 px-4 py-2 bg-fire-500 text-white rounded-lg hover:bg-fire-600 transition-colors flex items-center justify-center gap-2"
              aria-label="再生成"
            >
              <span>↻</span>
              <span>再生成</span>
            </button>
          )}
          {onDownload && (
            <button
              onClick={onDownload}
              className="flex-1 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors flex items-center justify-center gap-2"
              aria-label="ダウンロード"
            >
              <span>⬇</span>
              <span>ダウンロード</span>
            </button>
          )}
        </div>
      )}
    </div>
  )
}

export default ImagePreview
