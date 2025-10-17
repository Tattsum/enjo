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

/**
 * Mutation to post a tweet to Twitter
 */
export const POST_TO_TWITTER = gql`
  mutation PostToTwitter($input: TwitterPostInput!) {
    postToTwitter(input: $input) {
      success
      tweetId
      tweetUrl
      errorMessage
    }
  }
`

export interface TwitterPostInput {
  text: string
  imageUrl?: string
  addHashtag?: boolean
  addDisclaimer?: boolean
}

export interface TwitterPostResult {
  success: boolean
  tweetId?: string
  tweetUrl?: string
  errorMessage?: string
}

export interface PostToTwitterData {
  postToTwitter: TwitterPostResult
}

export interface PostToTwitterVariables {
  input: TwitterPostInput
}

/**
 * Mutation to generate an image from inflammatory text
 */
export const GENERATE_IMAGE = gql`
  mutation GenerateImage($input: GenerateImageInput!) {
    generateImage(input: $input) {
      imageUrl
      prompt
      generatedAt
    }
  }
`

export interface GenerateImageInput {
  text: string
  originalText?: string // Optional: original text before inflammatory conversion
  style?: ImageStyle
  aspectRatio?: AspectRatio
}

export enum ImageStyle {
  REALISTIC = 'REALISTIC',
  ILLUSTRATION = 'ILLUSTRATION',
  MEME = 'MEME',
  DRAMATIC = 'DRAMATIC',
}

export enum AspectRatio {
  SQUARE = 'SQUARE',
  LANDSCAPE = 'LANDSCAPE',
  PORTRAIT = 'PORTRAIT',
}

export interface GenerateImageResult {
  imageUrl: string
  prompt: string
  generatedAt: string
}

export interface GenerateImageData {
  generateImage: GenerateImageResult
}

export interface GenerateImageVariables {
  input: GenerateImageInput
}
