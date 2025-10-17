import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'

/**
 * Create Apollo Client instance for GraphQL communication
 * @returns Configured Apollo Client
 */
export function createApolloClient(): ApolloClient<unknown> {
  const httpLink = new HttpLink({
    uri: process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT || 'http://localhost:8080/graphql',
  })

  return new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'cache-and-network',
        errorPolicy: 'all',
      },
      query: {
        fetchPolicy: 'network-only',
        errorPolicy: 'all',
      },
      mutate: {
        errorPolicy: 'all',
      },
    },
  })
}
