/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      transitionProperty: {
        'max-height': 'max-height',
        'hide': 'color, background-color, border-color, opacity',
      },
      colors: {
        gray: {
          25: '#363b49',
          93: '#eaeef1',
          90: '#f0f3f5',
          95: '#f0f3f5',
          85: '#d1d9e0',
        },
        red: {
          55: '#e84f30',
        },
        purple: {
          55: '#271577',
          60: '#2f14a6',
        },
        blue: {
          65: '#5347d1',
          60: '#525ee1',
          55: '#5347d1',
        }
      },
    },
  },
  plugins: [],
};
