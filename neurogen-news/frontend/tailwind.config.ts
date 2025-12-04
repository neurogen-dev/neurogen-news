import type { Config } from 'tailwindcss'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Empatra Brand Colors - Authentic (Teal/Blue/Green palette)
        primary: {
          DEFAULT: '#14B8A6', // Teal - основной цвет Empatra
          50: '#F0FDFA',
          100: '#CCFBF1',
          200: '#99F6E4',
          300: '#5EEAD4',
          400: '#2DD4BF',
          500: '#14B8A6',
          600: '#0D9488',
          700: '#0F766E',
          800: '#115E59',
          900: '#134E4A',
        },
        // Secondary (dark blue/cyan)
        secondary: {
          DEFAULT: '#0891B2', // Dark cyan
          light: '#06B6D4',
          dark: '#0E7490',
        },
        // Accent (warm orange for highlights)
        accent: {
          DEFAULT: '#F97316', // Orange glow
          light: '#FB923C',
          dark: '#EA580C',
        },
        // Green (for nature/empathy theme)
        nature: {
          DEFAULT: '#10B981', // Emerald green
          light: '#34D399',
          dark: '#059669',
        },
        // Semantic colors
        success: '#10B981',
        warning: '#FBBF24',
        error: '#F43F5E',
        info: '#06B6D4',
        
        // Text colors for dark theme
        'text-primary': '#FAFAFA',
        'text-secondary': '#B8B8B8',
        'text-tertiary': '#808080',
        
        // Background colors - Deep dark (like Empatra's starry sky)
        'bg-base': '#0A0A0A', // Almost black
        'bg-elevated': '#0F172A', // Dark slate blue
        'bg-surface': '#1E293B', // Slate 800
        'bg-hover': '#334155', // Slate 700
        
        // Legacy backgrounds (for compatibility)
        'background-primary': '#0A0A0A',
        'background-secondary': '#0F172A',
        'background-tertiary': '#1E293B',
        
        // Border
        border: '#1E293B',
        'border-subtle': '#0F172A',
        
        // Glass effects for dark theme
        glass: {
          DEFAULT: 'rgba(30, 41, 59, 0.7)',
          light: 'rgba(51, 65, 85, 0.6)',
          dark: 'rgba(15, 23, 42, 0.9)',
        },
        
        // Empatra glow colors - Teal/Orange
        glow: {
          primary: 'rgba(20, 184, 166, 0.3)', // Teal glow
          secondary: 'rgba(8, 145, 178, 0.25)', // Cyan glow
          accent: 'rgba(249, 115, 22, 0.3)', // Orange glow
          nature: 'rgba(16, 185, 129, 0.25)', // Green glow
        },
        
        // Level badges
        beginner: '#34D399',
        intermediate: '#FBBF24',
        advanced: '#F43F5E',
      },
      fontFamily: {
        sans: [
          'Inter',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'sans-serif',
        ],
        display: [
          'Space Grotesk',
          'Inter',
          'sans-serif',
        ],
        mono: [
          'JetBrains Mono',
          'Fira Code',
          'Monaco',
          'Consolas',
          'monospace',
        ],
      },
      fontSize: {
        'xs': ['0.75rem', { lineHeight: '1rem' }],
        'sm': ['0.875rem', { lineHeight: '1.25rem' }],
        'base': ['1rem', { lineHeight: '1.6' }],
        'lg': ['1.125rem', { lineHeight: '1.75rem' }],
        'xl': ['1.25rem', { lineHeight: '1.75rem' }],
        '2xl': ['1.5rem', { lineHeight: '2rem' }],
        '3xl': ['1.875rem', { lineHeight: '2.25rem' }],
        '4xl': ['2.25rem', { lineHeight: '2.5rem' }],
        '5xl': ['3rem', { lineHeight: '1.2' }],
        'display': ['4rem', { lineHeight: '1.1' }],
      },
      spacing: {
        '18': '4.5rem',
        '22': '5.5rem',
        '30': '7.5rem',
        '112': '28rem',
        '128': '32rem',
      },
      borderRadius: {
        'sm': '6px',
        'DEFAULT': '8px',
        'md': '10px',
        'lg': '12px',
        'xl': '16px',
        '2xl': '20px',
        '3xl': '24px',
      },
      boxShadow: {
        'raised': '0 2px 4px rgba(0, 0, 0, 0.4), 0 4px 12px rgba(0, 0, 0, 0.3)',
        'floating': '0 4px 12px rgba(0, 0, 0, 0.5), 0 12px 32px rgba(0, 0, 0, 0.4)',
        'elevated': '0 8px 24px rgba(0, 0, 0, 0.6), 0 24px 48px rgba(0, 0, 0, 0.5)',
        'glow-primary': '0 0 30px rgba(20, 184, 166, 0.3)',
        'glow-secondary': '0 0 25px rgba(8, 145, 178, 0.25)',
        'glow-accent': '0 0 25px rgba(249, 115, 22, 0.3)',
        'glow-success': '0 0 20px rgba(16, 185, 129, 0.3)',
        'glow-danger': '0 0 20px rgba(244, 63, 94, 0.3)',
        'inset-soft': 'inset 0 2px 4px rgba(0, 0, 0, 0.2)',
        'inset-light': 'inset 0 1px 0 rgba(255, 255, 255, 0.08)',
      },
      backdropBlur: {
        'xs': '4px',
        'sm': '8px',
        'DEFAULT': '16px',
        'lg': '24px',
        'xl': '40px',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'slide-down': 'slideDown 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'scale-in': 'scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'blur-in': 'blurIn 0.6s cubic-bezier(0.4, 0, 0.2, 1)',
        'glow-pulse': 'glowPulse 3s ease-in-out infinite',
        'float': 'float 8s ease-in-out infinite',
        'shimmer': 'shimmer 2.5s linear infinite',
        'gradient-shift': 'gradientShift 8s ease-in-out infinite',
        'grain-shift': 'grainShift 15s steps(12) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
        blurIn: {
          '0%': { opacity: '0', filter: 'blur(12px)', transform: 'scale(0.95)' },
          '100%': { opacity: '1', filter: 'blur(0)', transform: 'scale(1)' },
        },
        glowPulse: {
          '0%, 100%': { filter: 'drop-shadow(0 0 12px rgba(20, 184, 166, 0.3))' },
          '50%': { filter: 'drop-shadow(0 0 24px rgba(20, 184, 166, 0.5))' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0) rotate(0deg)' },
          '50%': { transform: 'translateY(-12px) rotate(2deg)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        gradientShift: {
          '0%, 100%': { backgroundPosition: '0% 50%' },
          '50%': { backgroundPosition: '100% 50%' },
        },
        grainShift: {
          '0%': { backgroundPosition: '0 0, 0 0' },
          '100%': { backgroundPosition: '100% 100%, 200px 200px' },
        },
      },
      transitionTimingFunction: {
        'spring': 'cubic-bezier(0.34, 1.56, 0.64, 1)',
        'smooth': 'cubic-bezier(0.4, 0, 0.2, 1)',
        'bounce': 'cubic-bezier(0.68, -0.55, 0.265, 1.55)',
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
        'empatra-gradient': 'linear-gradient(135deg, #14B8A6 0%, #0891B2 50%, #10B981 100%)',
        'empatra-gradient-warm': 'linear-gradient(135deg, #14B8A6 0%, #F97316 50%, #10B981 100%)',
      },
      typography: (theme: (arg0: string) => unknown) => ({
        DEFAULT: {
          css: {
            maxWidth: 'none',
            color: theme('colors.text-primary'),
            a: {
              color: theme('colors.primary.DEFAULT'),
              '&:hover': {
                color: theme('colors.primary.400'),
              },
            },
            'code::before': { content: '""' },
            'code::after': { content: '""' },
            code: {
              backgroundColor: theme('colors.bg-surface'),
              padding: '0.25rem 0.375rem',
              borderRadius: '0.375rem',
              fontWeight: '400',
            },
          },
        },
        invert: {
          css: {
            color: theme('colors.zinc.300'),
            a: {
              color: theme('colors.primary.400'),
            },
            'h1, h2, h3, h4': {
              color: theme('colors.white'),
            },
            code: {
              backgroundColor: theme('colors.bg-surface'),
            },
            blockquote: {
              color: theme('colors.zinc.400'),
              borderLeftColor: theme('colors.bg-surface'),
            },
          },
        },
      }),
    },
  },
  plugins: [],
} satisfies Config
