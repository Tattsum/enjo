import React, { useState } from 'react';
import TwitterPostButton from './TwitterPostButton';

interface ResultDisplayProps {
  result: {
    original: string;
    inflammatory: string;
    explanation?: string;
  };
}

const ResultDisplay: React.FC<ResultDisplayProps> = ({ result }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(result.inflammatory);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Failed to copy text:', error);
    }
  };

  return (
    <div className="w-full space-y-6">
      {/* Before/After Comparison */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Original */}
        <div className="bg-white border-2 border-gray-200 rounded-lg p-6 shadow-sm">
          <div className="flex items-center gap-2 mb-4">
            <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
              <span className="text-blue-600 text-lg">👤</span>
            </div>
            <div>
              <p className="font-semibold text-gray-800">元の投稿</p>
              <p className="text-xs text-gray-500">通常の表現</p>
            </div>
          </div>
          <p className="text-gray-700 whitespace-pre-wrap">{result.original}</p>
        </div>

        {/* Inflammatory */}
        <div className="bg-gradient-to-br from-fire-50 to-fire-100 border-2 border-fire-300 rounded-lg p-6 shadow-lg relative">
          <div className="flex items-center gap-2 mb-4">
            <div className="w-10 h-10 rounded-full bg-fire-200 flex items-center justify-center">
              <span className="text-fire-700 text-lg">🔥</span>
            </div>
            <div>
              <p className="font-semibold text-fire-800">炎上化後</p>
              <p className="text-xs text-fire-600">要注意な表現</p>
            </div>
          </div>
          <p className="text-gray-800 whitespace-pre-wrap mb-4">{result.inflammatory}</p>

          <div className="flex gap-2 mt-4">
            <button
              onClick={handleCopy}
              className="px-4 py-2 bg-fire-600 text-white rounded-lg hover:bg-fire-700 transition-colors flex items-center gap-2"
              aria-label="コピー"
            >
              {copied ? (
                <>
                  <span>✓</span>
                  <span>コピーしました！</span>
                </>
              ) : (
                <>
                  <span>📋</span>
                  <span>コピー</span>
                </>
              )}
            </button>

            <TwitterPostButton text={result.inflammatory} />
          </div>
        </div>
      </div>

      {/* Explanation */}
      {result.explanation && (
        <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded">
          <div className="flex">
            <div className="flex-shrink-0">
              <span className="text-yellow-600 text-xl">💡</span>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-semibold text-yellow-800">説明</h3>
              <p className="mt-2 text-sm text-yellow-700">{result.explanation}</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ResultDisplay;
