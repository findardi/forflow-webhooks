<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<section class="panel empty">
	<div class="grille" aria-hidden="true"></div>
	{#if data.applications.length}
		<h1>Select an application</h1>
		<p>Pick one from the list to see and add its endpoints.</p>
	{:else}
		<h1>No applications yet</h1>
		<p>
			An application groups the endpoints that receive your webhooks. Create the first one on the
			left.
		</p>
	{/if}
</section>

<style>
	.empty {
		display: grid;
		justify-items: start;
		gap: 8px;
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
		max-width: 52ch;
		color: var(--ink-2);
	}

	@media (max-width: 800px) {
		.empty {
			display: none;
		}
	}
</style>
