import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TextInput from '../TextInput';

describe('TextInput', () => {
  it('renders textarea with placeholder', () => {
    const onChange = jest.fn();
    render(<TextInput value="" onChange={onChange} />);

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...');
    expect(textarea).toBeInTheDocument();
  });

  it('displays the current value', () => {
    const onChange = jest.fn();
    render(<TextInput value="テスト投稿" onChange={onChange} />);

    const textarea = screen.getByDisplayValue('テスト投稿');
    expect(textarea).toBeInTheDocument();
  });

  it('calls onChange when text is entered', async () => {
    const user = userEvent.setup();
    const onChange = jest.fn();
    render(<TextInput value="" onChange={onChange} />);

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...');
    await user.type(textarea, 'Hello');

    expect(onChange).toHaveBeenCalled();
  });

  it('limits input to 500 characters', () => {
    const onChange = jest.fn();
    render(<TextInput value="" onChange={onChange} />);

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...') as HTMLTextAreaElement;
    expect(textarea.maxLength).toBe(500);
  });

  it('displays character count', () => {
    const onChange = jest.fn();
    render(<TextInput value="テスト" onChange={onChange} />);

    expect(screen.getByText('3 / 500')).toBeInTheDocument();
  });

  it('shows warning when approaching character limit', () => {
    const onChange = jest.fn();
    const longText = 'a'.repeat(480);
    render(<TextInput value={longText} onChange={onChange} />);

    const counter = screen.getByText('480 / 500');
    expect(counter).toHaveClass('text-fire-500');
  });

  it('renders with proper accessibility attributes', () => {
    const onChange = jest.fn();
    render(<TextInput value="" onChange={onChange} />);

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...');
    expect(textarea).toHaveAttribute('aria-label', 'テキスト入力');
  });
});
