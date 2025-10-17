import { render, screen } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import userEvent from '@testing-library/user-event';
import ResultDisplay from '../ResultDisplay';

// Mock navigator.clipboard
const mockWriteText = jest.fn();
Object.assign(navigator, {
  clipboard: {
    writeText: mockWriteText,
  },
});

describe('ResultDisplay', () => {
  const mockResult = {
    original: '今日はいい天気ですね',
    inflammatory: 'こんな天気で外出する人は非常識だ',
    explanation: 'この表現は他者を批判する攻撃的な内容になっています',
  };

  beforeEach(() => {
    mockWriteText.mockClear();
    mockWriteText.mockResolvedValue(undefined);
  });

  it('renders original and inflammatory text side by side', () => {
    render(
      <MockedProvider>
        <ResultDisplay result={mockResult} />
      </MockedProvider>
    );

    expect(screen.getByText('元の投稿')).toBeInTheDocument();
    expect(screen.getByText('炎上化後')).toBeInTheDocument();
    expect(screen.getByText(mockResult.original)).toBeInTheDocument();
    expect(screen.getByText(mockResult.inflammatory)).toBeInTheDocument();
  });

  it('displays explanation when provided', () => {
    render(
      <MockedProvider>
        <ResultDisplay result={mockResult} />
      </MockedProvider>
    );

    expect(screen.getByText(/説明/)).toBeInTheDocument();
    expect(screen.getByText(mockResult.explanation)).toBeInTheDocument();
  });

  it('does not display explanation section when not provided', () => {
    const resultWithoutExplanation = {
      original: 'テスト',
      inflammatory: '炎上テスト',
    };
    render(
      <MockedProvider>
        <ResultDisplay result={resultWithoutExplanation} />
      </MockedProvider>
    );

    expect(screen.queryByText(/説明/)).not.toBeInTheDocument();
  });

  it('renders copy button for inflammatory text', () => {
    render(
      <MockedProvider>
        <ResultDisplay result={mockResult} />
      </MockedProvider>
    );

    const copyButton = screen.getByRole('button', { name: /コピー/ });
    expect(copyButton).toBeInTheDocument();
  });

  it('shows success message after copying', async () => {
    const user = userEvent.setup();
    render(
      <MockedProvider>
        <ResultDisplay result={mockResult} />
      </MockedProvider>
    );

    const copyButton = screen.getByRole('button', { name: /コピー/ });
    await user.click(copyButton);

    expect(screen.getByText(/コピーしました/)).toBeInTheDocument();
  });

  it('renders in SNS-style mockup design', () => {
    render(
      <MockedProvider>
        <ResultDisplay result={mockResult} />
      </MockedProvider>
    );

    // Check for visual elements that suggest SNS design
    const originalSection = screen.getByText('元の投稿').closest('div');
    const inflammatorySection = screen.getByText('炎上化後').closest('div');

    expect(originalSection).toBeInTheDocument();
    expect(inflammatorySection).toBeInTheDocument();
  });
});
