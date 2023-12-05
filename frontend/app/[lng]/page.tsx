import Link from 'next/link';
import { useTranslation } from '../i18n/index';

export default async function Page({ params: { lng } }:{
    params: { lng: string }
}) {
  const { t } = await useTranslation(lng, "")
  return (
    <>
      <h1>{t('title')}</h1>
      <Link href={`/${lng}/second-page`}>
        second page
      </Link>
    </>
  )
}
