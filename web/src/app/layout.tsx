import type { ReactNode } from "react";
import "./globals.scss";

type RootLayoutProps = {
  children: ReactNode;
};

const RootLayout = ({ children }: RootLayoutProps) => children;

export default RootLayout;
