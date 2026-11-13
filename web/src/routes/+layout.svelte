<script lang="ts">
	import '../app.css';
	import favicon from '#lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<title>Forflow</title>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
	<link
		rel="stylesheet"
		href="https://fonts.googleapis.com/css2?family=Hanken+Grotesk:wght@400;500;600&family=JetBrains+Mono:wght@400&display=swap"
	/>
</svelte:head>

<header>
	<a class="brand" href="/applications">Forflow</a>
	<nav aria-label="Main">
		<a href="/applications" aria-current="page">Applications</a>
	</nav>
</header>

<main>
	{@render children()}
</main>

<style>
	header {
		display: flex;
		align-items: center;
		gap: 40px;
		height: 64px;
		padding: 0 40px;
		border-bottom: 1px solid var(--rule);
	}

	.brand {
		font-size: 18px;
		font-weight: 600;
		text-decoration: none;
	}

	nav a {
		color: var(--ink-2);
		text-decoration: none;
	}

	nav a[aria-current='page'] {
		color: var(--ink);
		font-weight: 600;
	}

	main {
		padding: 28px 40px 40px;
	}

	@media (max-width: 800px) {
		header {
			padding: 0 16px;
			gap: 24px;
		}

		main {
			padding: 16px;
		}
	}
</style>
