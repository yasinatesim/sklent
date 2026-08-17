import { notFound } from "next/navigation";

import { NextIntlClientProvider } from "next-intl";
import { getMessages, setRequestLocale } from "next-intl/server";
import type { ReactNode } from "react";

import CategoryBar from "@/shared/layout/CategoryBar";
import ClientBootstrap from "@/shared/layout/ClientBootstrap";
import Footer from "@/shared/layout/Footer";
import Header from "@/shared/layout/Header";
import Toast from "@/shared/ui/Toast";

import { fetchCategories } from "@/features/web/catalog/api/products";

import CouponModal from "@/features/web/checkout/components/CouponModal";

import { routing } from "@/shared/i18n/routing";

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
