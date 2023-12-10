import { dir } from 'i18next';
import { languages } from '@/app/i18n/settings';
import type { Metadata } from 'next';
import '@/app/globals.css';
import { Inter } from 'next/font/google'
import Shell from './components/Footer/Shell';
const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'XOIDC Frontend',
  description: 'xoidc frontend',
}

export async function generateStaticParams() {
  return languages.map((lng) => ({ lng }))
}

export default function RootLayout({
  children,
  params: {
    lng
  }
}: {
  children: React.ReactNode,
  params: {
    lng: string
  }
}) {
  return (
    <html lang={lng} dir={dir(lng)}>
      <head />
      <body className={inter.className}>
        <Shell lng={lng} >
          {children}
        </Shell>
      </body>
    </html>
  )
}

