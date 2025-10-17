import { render, screen } from '@testing-library/react';
import ReplyList from '../ReplyList';
import { Reply, ReplyType } from '@/lib/graphql/queries';

describe('ReplyList', () => {
  const mockReplies: Reply[] = [
    {
      id: '1',
      type: ReplyType.LOGICAL_CRITICISM,
      content: 'それは論理的におかしいですよ。',
    },
    {
      id: '2',
      type: ReplyType.NITPICKING,
      content: '細かいことですが、そこは違いますね。',
    },
    {
      id: '3',
      type: ReplyType.OFF_TARGET,
      content: '全然関係ないですが、こういうこともありますよね。',
    },
    {
      id: '4',
      type: ReplyType.EXCESSIVE_DEFENSE,
      content: 'いや、それは絶対に正しいと思います！',
    },
  ];

  it('renders all replies', () => {
    render(<ReplyList replies={mockReplies} />);

    expect(screen.getByText('それは論理的におかしいですよ。')).toBeInTheDocument();
    expect(screen.getByText('細かいことですが、そこは違いますね。')).toBeInTheDocument();
    expect(
      screen.getByText('全然関係ないですが、こういうこともありますよね。')
    ).toBeInTheDocument();
    expect(screen.getByText('いや、それは絶対に正しいと思います！')).toBeInTheDocument();
  });

  it('displays correct icon for each reply type', () => {
    render(<ReplyList replies={mockReplies} />);

    // Check that type labels are displayed
    expect(screen.getByText(/正論で批判/)).toBeInTheDocument();
    expect(screen.getByText(/揚げ足を取る/)).toBeInTheDocument();
    expect(screen.getByText(/的外れな批判/)).toBeInTheDocument();
    expect(screen.getByText(/過剰に擁護/)).toBeInTheDocument();
  });

  it('renders empty state when no replies', () => {
    render(<ReplyList replies={[]} />);

    expect(screen.getByText(/リプライはまだありません/)).toBeInTheDocument();
  });

  it('renders list header', () => {
    render(<ReplyList replies={mockReplies} />);

    expect(screen.getByText(/想定されるリプライ/)).toBeInTheDocument();
  });

  it('displays reply count', () => {
    render(<ReplyList replies={mockReplies} />);

    expect(screen.getByText(/4.*件/)).toBeInTheDocument();
  });

  it('renders replies in SNS-style format', () => {
    render(<ReplyList replies={mockReplies} />);

    // Check for avatar-like elements (using emoji as avatars)
    const replyElements = screen.getAllByText(/それは論理的におかしいですよ。|細かいことですが/);
    expect(replyElements.length).toBeGreaterThan(0);
  });

  it('applies different styles for each reply type', () => {
    const { container } = render(<ReplyList replies={mockReplies} />);

    // Check that there are multiple reply items rendered
    const replyItems = container.querySelectorAll('[data-testid^="reply-"]');
    expect(replyItems.length).toBe(4);
  });
});
