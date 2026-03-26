/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				'on-error-container': '#93000a',
				'on-background': '#191c1e',
				'on-secondary': '#ffffff',
				'surface-dim': '#d8dadc',
				'surface-container-low': '#f2f4f6',
				'on-tertiary-fixed': '#00201c',
				'surface-tint': '#0056d2',
				'on-secondary-container': '#5c2900',
				'surface-bright': '#f7f9fb',
				'inverse-on-surface': '#eff1f3',
				'secondary-container': '#fd7c00',
				'on-secondary-fixed': '#311300',
				'surface-container-highest': '#e0e3e5',
				'error-container': '#ffdad6',
				'surface-container-high': '#e6e8ea',
				primary: '#0040a1',
				'on-primary': '#ffffff',
				'outline-variant': '#c3c6d6',
				'primary-fixed': '#dae2ff',
				error: '#ba1a1a',
				'on-secondary-fixed-variant': '#733500',
				'on-tertiary-container': '#6feed9',
				'on-tertiary-fixed-variant': '#005047',
				'on-tertiary': '#ffffff',
				'tertiary-fixed': '#79f7e3',
				'surface-variant': '#e0e3e5',
				'on-surface': '#191c1e',
				tertiary: '#005147',
				'tertiary-fixed-dim': '#59dbc7',
				background: '#f7f9fb',
				'on-surface-variant': '#424654',
				'inverse-surface': '#2d3133',
				outline: '#737785',
				secondary: '#984800',
				'on-primary-fixed': '#001847',
				'inverse-primary': '#b2c5ff',
				'secondary-fixed-dim': '#ffb689',
				'primary-container': '#0056d2',
				'on-error': '#ffffff',
				'surface-container': '#eceef0',
				'on-primary-container': '#ccd8ff',
				'surface-container-lowest': '#ffffff',
				'on-primary-fixed-variant': '#0040a1',
				surface: '#f7f9fb',
				'tertiary-container': '#006b5f',
				'primary-fixed-dim': '#b2c5ff',
				'secondary-fixed': '#ffdbc8'
			},
			fontFamily: {
				headline: ['Manrope', 'sans-serif'],
				body: ['Inter', 'sans-serif'],
				label: ['Inter', 'sans-serif']
			},
			borderRadius: {
				DEFAULT: '0.25rem',
				lg: '0.5rem',
				xl: '0.75rem',
				'2xl': '1rem',
				'3xl': '1.5rem',
				full: '9999px'
			}
		}
	},
	plugins: []
};
