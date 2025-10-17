'use client'

import React, { useState } from 'react'
import { useMutation } from '@apollo/client'
import { POST_TO_TWITTER, PostToTwitterData, PostToTwitterVariables } from '@/lib/graphql/queries'

interface TwitterPostButtonProps {
  text: string
  imageUrl?: string
  addHashtag?: boolean
  addDisclaimer?: boolean
}

const TwitterPostButton: React.FC<TwitterPostButtonProps> = ({
  text,
  imageUrl,
  addHashtag = true,
  addDisclaimer = true,
}) => {
  const [showDialog, setShowDialog] = useState(false)

  const [postToTwitter, { loading }] = useMutation<PostToTwitterData, PostToTwitterVariables>(
    POST_TO_TWITTER,
    {
      onCompleted: (data) => {
        if (data.postToTwitter.success) {
          alert('Twitterに投稿しました！')
          if (data.postToTwitter.tweetUrl) {
            window.open(data.postToTwitter.tweetUrl, '_blank')
          }
        } else {
          alert(`エラー: ${data.postToTwitter.errorMessage}`)
        }
        setShowDialog(false)
      },
      onError: (error) => {
        alert(`エラー: ${error.message}`)
        setShowDialog(false)
      },
    }
  )

  const handlePost = () => {
    postToTwitter({
      variables: {
        input: { text, imageUrl, addHashtag, addDisclaimer },
      },
    })
  }

  return (
    <>
      <button
        onClick={() => setShowDialog(true)}
        className="px-4 py-2 bg-black text-white rounded-lg hover:bg-gray-800 transition-colors flex items-center gap-2"
        disabled={loading}
        aria-label="𝕏に投稿"
      >
        {loading ? (
          <>
            <span>⏳</span>
            <span>投稿中...</span>
          </>
        ) : (
          <>
            <span>𝕏</span>
            <span>𝕏に投稿</span>
          </>
        )}
      </button>

      {showDialog && (
        <ConfirmDialog
          text={text}
          onConfirm={handlePost}
          onCancel={() => setShowDialog(false)}
          loading={loading}
        />
      )}
    </>
  )
}

interface ConfirmDialogProps {
  text: string
  onConfirm: () => void
  onCancel: () => void
  loading: boolean
}

const ConfirmDialog: React.FC<ConfirmDialogProps> = ({ text, onConfirm, onCancel, loading }) => {
  return (
    <div
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      onClick={onCancel}
    >
      <div
        className="bg-white rounded-lg p-6 max-w-md w-full mx-4 shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className="text-lg font-semibold text-gray-800 mb-4">
          ⚠️ Twitter/𝕏に投稿しますか？
        </h3>

        <div className="mb-4">
          <p className="text-sm text-gray-600 mb-2">投稿内容:</p>
          <div className="bg-gray-50 border border-gray-200 rounded p-3">
            <p className="text-gray-800 whitespace-pre-wrap">{text}</p>
          </div>
        </div>

        <div className="bg-yellow-50 border-l-4 border-yellow-400 p-3 mb-4">
          <p className="text-xs text-yellow-800">
            注意: この投稿は炎上シミュレーターで生成されたものです。
            <br />
            投稿による影響は自己責任となります。
          </p>
        </div>

        <div className="flex gap-3 justify-end">
          <button
            onClick={onCancel}
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition-colors"
            disabled={loading}
          >
            キャンセル
          </button>
          <button
            onClick={onConfirm}
            className="px-4 py-2 bg-black text-white rounded hover:bg-gray-800 transition-colors"
            disabled={loading}
          >
            {loading ? '投稿中...' : '投稿する'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default TwitterPostButton
