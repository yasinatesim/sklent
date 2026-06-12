import type { ReactNode } from "react";
import { notFound } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, setRequestLocale } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { fetchCategories } from "@/lib/api";
import ClientBootstrap from "@/components/layout/ClientBootstrap";
import Header from "@/components/layout/Header";
import CategoryBar from "@/components/layout/CategoryBar";
import Footer from "@/components/layout/Footer";
import Toast from "@/components/Toast/Toast";
import CouponModal from "@/components/CouponModal/CouponModal";

type LocaleLayoutProps = {
  children: ReactNode;
  params: Promise<{ locale: string }>;
};

export const generateStaticParams = () => routing.locales.map((locale) => ({ locale }));

const loadCategoriesSafely = async () => {
  try {
    return await fetchCategories();
  } catch {
    return [];
  }
};

const LocaleLayout = async ({ children, params }: LocaleLayoutProps) => {
  const { locale } = await params;
  if (!routing.locales.includes(locale as never)) {
    notFound();
  }
  setRequestLocale(locale);
  const messages = await getMessages();
  const categories = await loadCategoriesSafely();

  return (
    <html lang={locale}>
      <body>
        <NextIntlClientProvider locale={locale} messages={messages}>
          <ClientBootstrap />
          <Header />
          <CategoryBar categories={categories} />
          {children}
          <Footer />
          <Toast />
          <CouponModal />
        </NextIntlClientProvider>
      </body>
    </html>
  );
};

export default LocaleLayout;
