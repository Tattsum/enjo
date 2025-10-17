'use client'

import { ApolloProvider } from '@apollo/client'
import { createApolloClient } from '@/lib/graphql/client'
import './globals.css'

const apolloClient = createApolloClient()

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>
        <ApolloProvider client={apolloClient}>{children}</ApolloProvider>
      </body>
    </html>
  )
}
