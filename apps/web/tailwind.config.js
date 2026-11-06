/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // MOA premium identity: deep pine + champagne gold
        brand: {
          50: '#EBF2EC',
          100: '#D9E7DD',
          400: '#4E9068',
          500: '#2A7248',
          600: '#1E5A38', // Primary
          700: '#17452C', // Dark
          800: '#123522', // Pressed
          900: '#0C2318',
          950: '#08170F',
        },
        accent: {
          100: '#F2EAD3',
          500: '#C0A04E',
          600: '#9A7A33',
        },
        paper: '#F6F8F4',
        mist: '#ECEFE8',
        ink: '#141F19',
      },
      fontFamily: {
        sans: ['Pretendard', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Malgun Gothic', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
