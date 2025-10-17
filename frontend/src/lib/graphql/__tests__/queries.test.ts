import { GENERATE_INFLAMMATORY_TEXT, GENERATE_REPLIES } from '../queries'
import { gql } from '@apollo/client'

describe('GraphQL Queries', () => {
  describe('GENERATE_INFLAMMATORY_TEXT', () => {
    it('should be a valid GraphQL mutation', () => {
      expect(GENERATE_INFLAMMATORY_TEXT).toBeDefined()
      expect(GENERATE_INFLAMMATORY_TEXT.kind).toBe('Document')
    })

    it('should have correct mutation structure', () => {
      const expectedMutation = gql`
        mutation GenerateInflammatoryText($input: GenerateInput!) {
          generateInflammatoryText(input: $input) {
            inflammatoryText
            explanation
          }
        }
      `
      expect(GENERATE_INFLAMMATORY_TEXT.loc?.source.body.replace(/\s+/g, ' ').trim())
        .toBe(expectedMutation.loc?.source.body.replace(/\s+/g, ' ').trim())
    })
  })

  describe('GENERATE_REPLIES', () => {
    it('should be a valid GraphQL mutation', () => {
      expect(GENERATE_REPLIES).toBeDefined()
      expect(GENERATE_REPLIES.kind).toBe('Document')
    })

    it('should have correct mutation structure', () => {
      const expectedMutation = gql`
        mutation GenerateReplies($text: String!) {
          generateReplies(text: $text) {
            id
            type
            content
          }
        }
      `
      expect(GENERATE_REPLIES.loc?.source.body.replace(/\s+/g, ' ').trim())
        .toBe(expectedMutation.loc?.source.body.replace(/\s+/g, ' ').trim())
    })
  })
})
