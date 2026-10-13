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

export interface ActionOption {
  id: string;
  label: string;
  shortcut: string;
  run: () => Promise<void> | void;
}

export type ThemeType = 'default' | 'theme-catppuccin' | 'theme-tokyo' | 'theme-nord';

export type TabType = 'all' | 'apps' | 'clipboard';
