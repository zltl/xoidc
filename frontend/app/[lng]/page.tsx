import Link from "next/link";
import { useTranslation } from "@/app/i18n";
import { currentPath } from "../lib/server/currentPath";

export default async function Page({
  params: { lng },
}: {
  params: { lng: string };
}) {
  const { t } = await useTranslation(lng);
  const curpath = currentPath();

  console.log("current-path=", curpath);

  return (
    <div>
      this is main
      <div>
        <h1>{t("title")}</h1>
        <Link className='font-medium text-blue-600 dark:text-blue-500 hover:underline' href={`/${lng}/second-page`}>{t("to-second-page")}</Link>
        <br />
        <Link className='font-medium text-blue-600 dark:text-blue-500 hover:underline' href={`/${lng}/client-page`}>{t("to-client-page")}</Link>
        <br />
        <div>{t("this is a test text")}</div>
      </div>
      <div className="h-[3000px]">abcdefg</div>
    </div>
  );
}
