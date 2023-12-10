import Link from 'next/link'
import { Trans } from 'react-i18next/TransWithoutContext'
import { languages } from '@/app/i18n/settings'
import { useCurrentPath } from '@/app/lib/client/useCurrentPath'
import { currentPath } from '@/app/lib/server/currentPath'

export const FooterBase = ({ t, lng }: { t: any, lng: string }) => {

  const curpath = currentPath();

  return (
    <footer style={{ marginTop: 50 }}>
      <Trans i18nKey="languageSwitcher" t={t}>
        Switch from <strong>{{ lng }}</strong> to:{' '}
      </Trans>
      {languages.filter((l) => lng !== l).map((l, index) => {
        return (
          <span key={l}>
            {index > 0 && (' or ')}
            <Link href={`/${l}/${curpath}`}>
              {l}
            </Link>
          </span>
        )
      })}
    </footer>
  )
}
