import Link from 'next/link';
import { useTranslation } from '@/app/i18n';
import { Footer } from '@/app/[lng]/components/Footer'

export default async function Page({ params: { lng } }: {
  params: { lng: string }
}) {
  const { t } = await useTranslation(lng, 'second-page')

  return (
    <>
      <h1>{t('title')}</h1>
      <Link href={`/${lng}`}>
        {t('back-to-home')}
      </Link>
      <Footer lng={lng} />
    </>
  )
}
