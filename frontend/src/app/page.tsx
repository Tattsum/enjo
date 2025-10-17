'use client'

import React, { useState } from 'react'
import { useMutation } from '@apollo/client'
import TextInput from '@/components/TextInput'
import LevelSlider from '@/components/LevelSlider'
import ResultDisplay from '@/components/ResultDisplay'
import ReplyList from '@/components/ReplyList'
import {
  GENERATE_INFLAMMATORY_TEXT,
  GENERATE_REPLIES,
  GenerateInflammatoryTextData,
  GenerateInflammatoryTextVariables,
  GenerateRepliesData,
  GenerateRepliesVariables,
  Reply,
} from '@/lib/graphql/queries'

export default function Home() {
  // 状態管理
  const [inputText, setInputText] = useState<string>('')
  const [level, setLevel] = useState<number>(3)
  const [result, setResult] = useState<{
    original: string
    inflammatory: string
    explanation?: string
  } | null>(null)
  const [replies, setReplies] = useState<Reply[]>([])
  const [errorMessage, setErrorMessage] = useState<string>('')

  // GraphQL Mutation
  const [generateInflammatoryText, { loading: generateLoading }] = useMutation<
    GenerateInflammatoryTextData,
    GenerateInflammatoryTextVariables
  >(GENERATE_INFLAMMATORY_TEXT, {
    onCompleted: (data) => {
      setResult({
        original: inputText,
        inflammatory: data.generateInflammatoryText.inflammatoryText,
        explanation: data.generateInflammatoryText.explanation,
      })
      setReplies([])
      setErrorMessage('')
    },
    onError: (error) => {
      console.error('Error generating inflammatory text:', error)
      setErrorMessage('エラーが発生しました: ' + error.message)
    },
  })

  const [generateReplies, { loading: repliesLoading }] = useMutation<
    GenerateRepliesData,
    GenerateRepliesVariables
  >(GENERATE_REPLIES, {
    onCompleted: (data) => {
      setReplies(data.generateReplies)
      setErrorMessage('')
    },
    onError: (error) => {
      console.error('Error generating replies:', error)
      setErrorMessage('エラーが発生しました: ' + error.message)
    },
  })

  // ハンドラー
  const handleGenerate = async () => {
    if (!inputText.trim()) return

    await generateInflammatoryText({
      variables: {
        input: {
          originalText: inputText,
          level,
        },
      },
    })
  }

  const handleGenerateReplies = async () => {
    if (!result) return

    await generateReplies({
      variables: {
        text: result.inflammatory,
      },
    })
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <header className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-fire-600 mb-4">
            🔥 炎上シミュレーター
          </h1>
          <p className="text-lg text-gray-700 mb-2">
            SNS投稿を「炎上しやすい表現」に変換して、リスクを学ぶツール
          </p>
          <p className="text-sm text-yellow-700 bg-yellow-50 border border-yellow-200 rounded-lg px-4 py-2 inline-block">
            ⚠️ このツールは教育・エンターテインメント目的です
          </p>
        </header>

        {/* Input Section */}
        <div className="bg-white rounded-xl shadow-lg p-6 md:p-8 mb-8">
          <div className="space-y-6">
            <div>
              <h2 className="text-xl font-semibold text-gray-800 mb-4">
                テキストを入力
              </h2>
              <TextInput value={inputText} onChange={setInputText} />
            </div>

            <div>
              <LevelSlider value={level} onChange={setLevel} />
            </div>

            <button
              onClick={handleGenerate}
              disabled={!inputText.trim() || generateLoading}
              className="w-full py-4 px-6 bg-fire-600 text-white rounded-lg font-semibold text-lg hover:bg-fire-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
              aria-label="🔥 炎上化する"
            >
              {generateLoading ? (
                <>
                  <span className="animate-spin">⏳</span>
                  <span>生成中...</span>
                </>
              ) : (
                <>
                  <span>🔥</span>
                  <span>炎上化する</span>
                </>
              )}
            </button>
          </div>
        </div>

        {/* Error Message */}
        {errorMessage && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-8">
            {errorMessage}
          </div>
        )}

        {/* Result Section */}
        {result && (
          <div className="bg-white rounded-xl shadow-lg p-6 md:p-8 mb-8">
            <ResultDisplay result={result} />

            <div className="mt-8 pt-6 border-t border-gray-200">
              <button
                onClick={handleGenerateReplies}
                disabled={repliesLoading}
                className="w-full py-3 px-6 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
                aria-label="💬 リプライを生成"
              >
                {repliesLoading ? (
                  <>
                    <span className="animate-spin">⏳</span>
                    <span>生成中...</span>
                  </>
                ) : (
                  <>
                    <span>💬</span>
                    <span>リプライを生成</span>
                  </>
                )}
              </button>
            </div>
          </div>
        )}

        {/* Replies Section */}
        {replies.length > 0 && (
          <div className="bg-white rounded-xl shadow-lg p-6 md:p-8">
            <ReplyList replies={replies} />
          </div>
        )}
      </div>
    </main>
  )
}
