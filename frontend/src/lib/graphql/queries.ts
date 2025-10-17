import { gql } from '@apollo/client'

/**
 * Mutation to generate inflammatory text from original text
 */
export const GENERATE_INFLAMMATORY_TEXT = gql`
  mutation GenerateInflammatoryText($input: GenerateInput!) {
    generateInflammatoryText(input: $input) {
      inflammatoryText
      explanation
    }
  }
`

/**
 * Mutation to generate replies for a given text
 */
export const GENERATE_REPLIES = gql`
  mutation GenerateReplies($text: String!) {
    generateReplies(text: $text) {
      id
      type
      content
    }
  }
`

/**
 * Type definitions for TypeScript
 */

export interface GenerateInput {
  originalText: string
  level: number
}

export interface GenerateResult {
  inflammatoryText: string
  explanation?: string
}

export interface Reply {
  id: string
  type: ReplyType
  content: string
}

export enum ReplyType {
  LOGICAL_CRITICISM = 'LOGICAL_CRITICISM',
  NITPICKING = 'NITPICKING',
  OFF_TARGET = 'OFF_TARGET',
  EXCESSIVE_DEFENSE = 'EXCESSIVE_DEFENSE',
}

export interface GenerateInflammatoryTextData {
  generateInflammatoryText: GenerateResult
}

export interface GenerateInflammatoryTextVariables {
  input: GenerateInput
}

export interface GenerateRepliesData {
  generateReplies: Reply[]
}

export interface GenerateRepliesVariables {
  text: string
}
