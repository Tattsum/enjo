import React from 'react'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import Home from '../page'
import {
  GENERATE_INFLAMMATORY_TEXT,
  GENERATE_REPLIES,
  ReplyType,
} from '@/lib/graphql/queries'

describe('Home Page', () => {
  it('renders the header and title', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    expect(screen.getByText('🔥 炎上シミュレーター')).toBeInTheDocument()
    expect(
      screen.getByText('⚠️ このツールは教育・エンターテインメント目的です')
    ).toBeInTheDocument()
  })

  it('renders TextInput and LevelSlider', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    expect(screen.getByPlaceholderText('普通の投稿を入力してください...')).toBeInTheDocument()
    expect(screen.getByLabelText('炎上レベル')).toBeInTheDocument()
  })

  it('renders the generate button', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    expect(screen.getByRole('button', { name: /🔥 炎上化する/ })).toBeInTheDocument()
  })

  it('disables the generate button when input is empty', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    expect(button).toBeDisabled()
  })

  it('enables the generate button when input has text', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    expect(button).not.toBeDisabled()
  })

  it('shows loading state when generating inflammatory text', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_INFLAMMATORY_TEXT,
          variables: {
            input: {
              originalText: 'テスト投稿',
              level: 3,
            },
          },
        },
        result: {
          data: {
            generateInflammatoryText: {
              inflammatoryText: '炎上化されたテキスト',
              explanation: '説明文',
            },
          },
        },
        delay: 100,
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    fireEvent.click(button)

    expect(screen.getByText(/生成中/)).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.queryByText(/生成中/)).not.toBeInTheDocument()
    })
  })

  it('generates and displays inflammatory text', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_INFLAMMATORY_TEXT,
          variables: {
            input: {
              originalText: 'テスト投稿',
              level: 3,
            },
          },
        },
        result: {
          data: {
            generateInflammatoryText: {
              inflammatoryText: '炎上化されたテキスト',
              explanation: '説明文',
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    fireEvent.click(button)

    await waitFor(() => {
      expect(screen.getByText('元の投稿')).toBeInTheDocument()
      expect(screen.getByText('炎上化後')).toBeInTheDocument()
      expect(screen.getByText('炎上化されたテキスト')).toBeInTheDocument()
    })
  })

  it('shows generate replies button after generating inflammatory text', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_INFLAMMATORY_TEXT,
          variables: {
            input: {
              originalText: 'テスト投稿',
              level: 3,
            },
          },
        },
        result: {
          data: {
            generateInflammatoryText: {
              inflammatoryText: '炎上化されたテキスト',
              explanation: '説明文',
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    fireEvent.click(button)

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /💬 リプライを生成/ })
      ).toBeInTheDocument()
    })
  })

  it('generates and displays replies', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_INFLAMMATORY_TEXT,
          variables: {
            input: {
              originalText: 'テスト投稿',
              level: 3,
            },
          },
        },
        result: {
          data: {
            generateInflammatoryText: {
              inflammatoryText: '炎上化されたテキスト',
              explanation: '説明文',
            },
          },
        },
      },
      {
        request: {
          query: GENERATE_REPLIES,
          variables: {
            text: '炎上化されたテキスト',
          },
        },
        result: {
          data: {
            generateReplies: [
              {
                id: '1',
                type: ReplyType.LOGICAL_CRITICISM,
                content: '正論で批判するコメント',
              },
              {
                id: '2',
                type: ReplyType.NITPICKING,
                content: '揚げ足を取るコメント',
              },
              {
                id: '3',
                type: ReplyType.OFF_TARGET,
                content: '的外れな批判コメント',
              },
              {
                id: '4',
                type: ReplyType.EXCESSIVE_DEFENSE,
                content: '過剰に擁護するコメント',
              },
            ],
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const generateButton = screen.getByRole('button', { name: /🔥 炎上化する/ })
    fireEvent.click(generateButton)

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /💬 リプライを生成/ })).toBeInTheDocument()
    })

    const replyButton = screen.getByRole('button', { name: /💬 リプライを生成/ })
    fireEvent.click(replyButton)

    await waitFor(() => {
      expect(screen.getByText('正論で批判するコメント')).toBeInTheDocument()
      expect(screen.getByText('揚げ足を取るコメント')).toBeInTheDocument()
      expect(screen.getByText('的外れな批判コメント')).toBeInTheDocument()
      expect(screen.getByText('過剰に擁護するコメント')).toBeInTheDocument()
    })
  })

  it('handles error when generating inflammatory text fails', async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_INFLAMMATORY_TEXT,
          variables: {
            input: {
              originalText: 'テスト投稿',
              level: 3,
            },
          },
        },
        error: new Error('APIエラー'),
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const textarea = screen.getByPlaceholderText('普通の投稿を入力してください...')
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } })

    const button = screen.getByRole('button', { name: /🔥 炎上化する/ })
    fireEvent.click(button)

    await waitFor(() => {
      expect(screen.getByText(/エラーが発生しました/)).toBeInTheDocument()
    })
  })

  it('updates level when slider changes', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <Home />
      </MockedProvider>
    )

    const slider = screen.getByLabelText('炎上レベル')
    fireEvent.change(slider, { target: { value: '5' } })

    expect(screen.getByText('レベル 5')).toBeInTheDocument()
    expect(screen.getByText('炎上確実な表現')).toBeInTheDocument()
  })
})
