<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>{page.status} · Forflow</title>
</svelte:head>

<section class="panel">
	<div class="grille" aria-hidden="true"></div>
	<h1>{page.status === 404 ? 'Not found' : 'Something went wrong'}</h1>
	<p>{page.error?.message}</p>
	<a href="/applications">Back to applications</a>
</section>

<style>
	section {
		display: grid;
		justify-items: start;
		gap: 8px;
		max-width: 640px;
		margin: 40px auto;
		padding: 40px 28px;
	}

	.grille {
		width: 96px;
		height: 48px;
		margin-bottom: 16px;
	}

	h1 {
		font-size: 20px;
		font-weight: 600;
	}

	p {
		color: var(--ink-2);
	}
</style>
