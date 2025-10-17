import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import TwitterPostButton from '../TwitterPostButton'
import { POST_TO_TWITTER } from '@/lib/graphql/queries'

// Mock window.open
const mockWindowOpen = jest.fn()
window.open = mockWindowOpen

// Mock window.alert
const mockAlert = jest.fn()
window.alert = mockAlert

describe('TwitterPostButton', () => {
  beforeEach(() => {
    mockWindowOpen.mockClear()
    mockAlert.mockClear()
  })

  it('renders post button', () => {
    render(
      <MockedProvider>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )
    expect(screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })).toBeInTheDocument()
  })

  it('shows confirmation dialog on click', () => {
    render(
      <MockedProvider>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    expect(screen.getByText(/投稿しますか/i)).toBeInTheDocument()
    expect(screen.getByText('テストツイート')).toBeInTheDocument()
  })

  it('closes dialog on cancel', () => {
    render(
      <MockedProvider>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    const cancelButton = screen.getByRole('button', { name: /キャンセル/i })
    fireEvent.click(cancelButton)

    expect(screen.queryByText(/投稿しますか/i)).not.toBeInTheDocument()
  })

  it('posts to Twitter on confirmation', async () => {
    const mocks = [
      {
        request: {
          query: POST_TO_TWITTER,
          variables: {
            input: {
              text: 'テストツイート',
              addHashtag: true,
              addDisclaimer: true,
            },
          },
        },
        result: {
          data: {
            postToTwitter: {
              success: true,
              tweetId: '123456789',
              tweetUrl: 'https://twitter.com/user/status/123456789',
              errorMessage: null,
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    const confirmButton = screen.getByRole('button', { name: /投稿する/i })
    fireEvent.click(confirmButton)

    await waitFor(() => {
      expect(mockAlert).toHaveBeenCalledWith('Twitterに投稿しました！')
    })

    await waitFor(() => {
      expect(mockWindowOpen).toHaveBeenCalledWith('https://twitter.com/user/status/123456789', '_blank')
    })
  })

  it('shows error message on post failure', async () => {
    const mocks = [
      {
        request: {
          query: POST_TO_TWITTER,
          variables: {
            input: {
              text: 'テストツイート',
              addHashtag: true,
              addDisclaimer: true,
            },
          },
        },
        result: {
          data: {
            postToTwitter: {
              success: false,
              tweetId: null,
              tweetUrl: null,
              errorMessage: '投稿に失敗しました',
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    const confirmButton = screen.getByRole('button', { name: /投稿する/i })
    fireEvent.click(confirmButton)

    await waitFor(() => {
      expect(mockAlert).toHaveBeenCalledWith('エラー: 投稿に失敗しました')
    })
  })

  it('shows loading state while posting', async () => {
    const mocks = [
      {
        request: {
          query: POST_TO_TWITTER,
          variables: {
            input: {
              text: 'テストツイート',
              addHashtag: true,
              addDisclaimer: true,
            },
          },
        },
        result: {
          data: {
            postToTwitter: {
              success: true,
              tweetId: '123456789',
              tweetUrl: 'https://twitter.com/user/status/123456789',
              errorMessage: null,
            },
          },
        },
        delay: 100,
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <TwitterPostButton text="テストツイート" />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    const confirmButton = screen.getByRole('button', { name: /投稿する/i })

    // Confirm button should be clickable before loading
    expect(confirmButton).not.toBeDisabled()

    fireEvent.click(confirmButton)

    // Dialog should show "投稿中..." on the confirm button
    await waitFor(() => {
      const loadingButton = screen.getByRole('button', { name: /投稿中/i })
      expect(loadingButton).toBeInTheDocument()
      expect(loadingButton).toBeDisabled()
    })
  })

  it('supports custom hashtag and disclaimer options', async () => {
    const mocks = [
      {
        request: {
          query: POST_TO_TWITTER,
          variables: {
            input: {
              text: 'テストツイート',
              addHashtag: false,
              addDisclaimer: false,
            },
          },
        },
        result: {
          data: {
            postToTwitter: {
              success: true,
              tweetId: '123456789',
              tweetUrl: 'https://twitter.com/user/status/123456789',
              errorMessage: null,
            },
          },
        },
      },
    ]

    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <TwitterPostButton text="テストツイート" addHashtag={false} addDisclaimer={false} />
      </MockedProvider>
    )

    const button = screen.getByRole('button', { name: /Twitterに投稿|𝕏に投稿/i })
    fireEvent.click(button)

    const confirmButton = screen.getByRole('button', { name: /投稿する/i })
    fireEvent.click(confirmButton)

    await waitFor(() => {
      expect(mockAlert).toHaveBeenCalledWith('Twitterに投稿しました！')
    })
  })
})
