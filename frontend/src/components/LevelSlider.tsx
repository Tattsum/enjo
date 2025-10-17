import React from 'react';

interface LevelSliderProps {
  value: number;
  onChange: (value: number) => void;
}

const levelDescriptions: Record<number, string> = {
  1: '少し配慮に欠ける表現',
  2: '誤解を招きやすい表現',
  3: '明確に批判されそうな表現',
  4: 'かなり問題がある表現',
  5: '炎上確実な表現',
};

const levelColors: Record<number, string> = {
  1: 'accent-blue-500',
  2: 'accent-yellow-500',
  3: 'accent-orange-500',
  4: 'accent-fire-600',
  5: 'accent-fire-700',
};

const LevelSlider: React.FC<LevelSliderProps> = ({ value, onChange }) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(parseInt(e.target.value, 10));
  };

  return (
    <div className="w-full">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-gray-700">炎上レベル</h3>
        <span className="text-2xl font-bold text-fire-600">レベル {value}</span>
      </div>

      <input
        type="range"
        min="1"
        max="5"
        value={value}
        onChange={handleChange}
        aria-label="炎上レベル"
        className={`w-full h-2 rounded-lg appearance-none cursor-pointer ${levelColors[value]}`}
      />

      <div className="flex justify-between text-xs text-gray-500 mt-2">
        <span>1</span>
        <span>2</span>
        <span>3</span>
        <span>4</span>
        <span>5</span>
      </div>

      <div className="mt-4 p-4 bg-gray-50 rounded-lg">
        <p className="text-sm text-gray-700">{levelDescriptions[value]}</p>
      </div>
    </div>
  );
};

export default LevelSlider;
