import React from 'react';

interface TextInputProps {
  value: string;
  onChange: (value: string) => void;
}

const TextInput: React.FC<TextInputProps> = ({ value, onChange }) => {
  const maxLength = 500;
  const currentLength = value.length;
  const isNearLimit = currentLength >= 450;

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onChange(e.target.value);
  };

  return (
    <div className="w-full">
      <textarea
        value={value}
        onChange={handleChange}
        placeholder="普通の投稿を入力してください..."
        maxLength={maxLength}
        aria-label="テキスト入力"
        className="w-full h-32 p-4 border-2 border-gray-300 rounded-lg focus:border-fire-500 focus:outline-none resize-none transition-colors text-gray-900 placeholder:text-gray-400"
        rows={5}
      />
      <div className="flex justify-end mt-2">
        <span
          className={`text-sm ${
            isNearLimit ? 'text-fire-500 font-bold' : 'text-gray-500'
          }`}
        >
          {currentLength} / {maxLength}
        </span>
      </div>
    </div>
  );
};

export default TextInput;
