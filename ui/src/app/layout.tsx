import type { Metadata } from "next";
import { Inter, Lexend } from 'next/font/google'
import "./globals.css";
import clsx from 'clsx'

export const metadata: Metadata = {
    title: {
      template: '%s - myLocal',
      default: 'myLocal - Authentication',
    },
    description:
      'myLocal - authentication for the myLocal app.',
  }
  
  const inter = Inter({
    subsets: ['latin'],
    display: 'swap',
    variable: '--font-inter',
  })
  
  const lexend = Lexend({
    subsets: ['latin'],
    display: 'swap',
    variable: '--font-lexend',
  })
  
  export default function RootLayout({
    children,
  }: {
    children: React.ReactNode
  }) {
    return (
      <html
        lang="en"
        className={clsx(
          'h-full scroll-smooth bg-slate-50 antialiased',
          inter.variable,
          lexend.variable,
        )}
      >
        <body className="flex h-full flex-col">{children}</body>
      </html>
    )
  }
