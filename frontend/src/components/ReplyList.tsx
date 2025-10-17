import React from 'react';
import { Reply, ReplyType } from '@/lib/graphql/queries';

interface ReplyListProps {
  replies: Reply[];
}

const replyTypeLabels: Record<ReplyType, string> = {
  [ReplyType.LOGICAL_CRITICISM]: '正論で批判',
  [ReplyType.NITPICKING]: '揚げ足を取る',
  [ReplyType.OFF_TARGET]: '的外れな批判',
  [ReplyType.EXCESSIVE_DEFENSE]: '過剰に擁護',
};

const replyTypeIcons: Record<ReplyType, string> = {
  [ReplyType.LOGICAL_CRITICISM]: '🤓',
  [ReplyType.NITPICKING]: '🔍',
  [ReplyType.OFF_TARGET]: '🎯',
  [ReplyType.EXCESSIVE_DEFENSE]: '🛡️',
};

const replyTypeColors: Record<ReplyType, string> = {
  [ReplyType.LOGICAL_CRITICISM]: 'bg-blue-50 border-blue-200',
  [ReplyType.NITPICKING]: 'bg-yellow-50 border-yellow-200',
  [ReplyType.OFF_TARGET]: 'bg-purple-50 border-purple-200',
  [ReplyType.EXCESSIVE_DEFENSE]: 'bg-green-50 border-green-200',
};

const ReplyList: React.FC<ReplyListProps> = ({ replies }) => {
  if (replies.length === 0) {
    return (
      <div className="w-full p-8 text-center">
        <p className="text-gray-500">リプライはまだありません</p>
      </div>
    );
  }

  return (
    <div className="w-full">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-xl font-semibold text-gray-800">想定されるリプライ</h3>
        <span className="text-sm text-gray-600">{replies.length} 件</span>
      </div>

      <div className="space-y-4">
        {replies.map((reply, index) => (
          <div
            key={reply.id}
            data-testid={`reply-${reply.id}`}
            className={`border-2 rounded-lg p-4 transition-all hover:shadow-md ${
              replyTypeColors[reply.type]
            } animate-fade-in`}
            style={{
              animationDelay: `${index * 100}ms`,
            }}
          >
            <div className="flex items-start gap-3">
              {/* Avatar */}
              <div className="flex-shrink-0 w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-xl">
                {replyTypeIcons[reply.type]}
              </div>

              {/* Content */}
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-2">
                  <span className="font-semibold text-gray-800">
                    {replyTypeLabels[reply.type]}
                  </span>
                  <span className="text-xs text-gray-500">タイプ</span>
                </div>
                <p className="text-gray-700 whitespace-pre-wrap">{reply.content}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <style jsx>{`
        @keyframes fade-in {
          from {
            opacity: 0;
            transform: translateY(10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .animate-fade-in {
          animation: fade-in 0.3s ease-out forwards;
          opacity: 0;
        }
      `}</style>
    </div>
  );
};

export default ReplyList;
