export interface AppItem {
  id: string;
  name: string;
  comment: string;
  exec: string;
  icon: string;
  iconPath: string;
  terminal: boolean;
  desktopPath: string;
  score: number;
  launchCount: number;
  resolvedIcon?: string;
}

export interface ClipboardEntry {
  id: string;
  content: string;
  preview: string;
  timestamp: string;
  charCount: number;
}

export interface CalcItem {
  expression: string;
  result: string;
  success: boolean;
}

export interface TodoEntry {
  id: string;
  title: string;
  completed: boolean;
  createdAt: string;
  completedAt?: string;
}

export interface EmojiEntry {
  emoji: string;
  name: string;
  shortcode: string;
  category: string;
  keywords: string[];
}

export interface ConvertItem {
  success: boolean;
  query: string;
  fromValue: string;
  fromUnit: string;
  toValue: string;
  toUnit: string;
  result: string;
  category: string;
  formatted: string;
}

export interface GuideShortcut {
  key: string;
  desc: string;
}

export interface GuideTopic {
  id: string;
  title: string;
  category: string;
  summary: string;
  content: string[];
  shortcuts: GuideShortcut[];
  examples: string[];
}

export interface ActionOption {
  id: string;
  label: string;
  shortcut: string;
  run: () => Promise<void> | void;
}

export type ThemeType = 'default' | 'theme-catppuccin' | 'theme-tokyo' | 'theme-nord';

export type TabType = 'all' | 'apps' | 'clipboard' | 'todo' | 'emoji' | 'guide';


