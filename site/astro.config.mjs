// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://cloverhound.github.io',
	base: process.env.BASE_URL || '/webex-cli',
	integrations: [
		starlight({
			title: 'Webex CLI',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/Cloverhound/webex-cli' }],
			sidebar: [
				{ label: 'Getting Started', items: [
					{ label: 'Installation', slug: 'installation' },
					{ label: 'Quick Start', slug: 'getting-started' },
				]},
				{ label: 'API Reference', items: [
					{ label: 'Admin Commands', slug: 'admin' },
					{ label: 'Calling Commands', slug: 'calling' },
					{ label: 'Contact Center Commands', slug: 'cc' },
					{ label: 'Device Commands', slug: 'device' },
					{ label: 'Meetings Commands', slug: 'meetings' },
					{ label: 'Messaging Commands', slug: 'messaging' },
				]},
				{ label: 'Integrations', items: [
					{ label: 'Claude Code Skill', slug: 'claude-skill' },
				]},
			],
		}),
	],
});
