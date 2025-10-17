import { createApolloClient } from '../client'
import { ApolloClient, InMemoryCache } from '@apollo/client'

// Mock fetch for Node.js environment
global.fetch = jest.fn()

describe('Apollo Client', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('should create Apollo Client instance', () => {
    const client = createApolloClient()
    expect(client).toBeInstanceOf(ApolloClient)
  })

  it('should have InMemoryCache', () => {
    const client = createApolloClient()
    expect(client.cache).toBeInstanceOf(InMemoryCache)
  })

  it('should use correct GraphQL endpoint from env', () => {
    const originalEnv = process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT
    process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT = 'http://test:8080/graphql'

    const client = createApolloClient()
    // @ts-expect-error - accessing private property for testing
    const uri = client.link?.options?.uri
    expect(uri).toBe('http://test:8080/graphql')

    process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT = originalEnv
  })

  it('should use default endpoint if env is not set', () => {
    const originalEnv = process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT
    delete process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT

    const client = createApolloClient()
    // @ts-expect-error - accessing private property for testing
    const uri = client.link?.options?.uri
    expect(uri).toBe('http://localhost:8080/graphql')

    process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT = originalEnv
  })
})
