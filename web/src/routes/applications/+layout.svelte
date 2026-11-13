<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	let creating = $state(false);
	const createError = $derived(
		page.form && 'name' in page.form ? (page.form.message as string) : undefined
	);
</script>

<div class="split" class:has-selection={page.params.id}>
	<section class="panel list" aria-labelledby="apps-heading">
		<div class="list-head">
			<h2 id="apps-heading" class="label">Applications</h2>
			<span class="label count">{data.applications.length}</span>
		</div>

		{#if data.applications.length}
			<ul>
				{#each data.applications as app (app.id)}
					<li>
						<a
							href="/applications/{app.id}"
							aria-current={page.params.id === app.id ? 'page' : undefined}
						>
							<span class="name">{app.name}</span>
							<span class="mono short">{app.id.slice(0, 8)}</span>
						</a>
					</li>
				{/each}
			</ul>
		{/if}

		<form
			method="POST"
			action="?/create"
			use:enhance={() => {
				creating = true;
				return async ({ update }) => {
					await update();
					creating = false;
				};
			}}
		>
			<label for="app-name" class="label">New application</label>
			<div class="row">
				<input
					id="app-name"
					class="field"
					name="name"
					required
					maxlength="100"
					autocomplete="off"
					placeholder="billing-service"
					aria-invalid={createError ? 'true' : undefined}
					aria-describedby={createError ? 'app-name-error' : undefined}
				/>
				<button
					class="btn"
					class:btn-key={!page.params.id}
					class:btn-ghost={!!page.params.id}
					disabled={creating}>Create</button
				>
			</div>
			{#if createError}
				<p id="app-name-error" class="error">{createError}</p>
			{/if}
		</form>
	</section>

	<div class="detail">
		{@render children()}
	</div>
</div>

<style>
	.split {
		display: grid;
		grid-template-columns: minmax(280px, 380px) minmax(0, 1fr);
		gap: 24px;
		align-items: start;
		max-width: 1600px;
		margin: 0 auto;
	}

	.list-head {
		display: flex;
		justify-content: space-between;
		padding: 18px 22px;
		border-bottom: 1px solid var(--rule);
	}

	.count {
		font-variant-numeric: tabular-nums;
	}

	ul {
		margin: 0;
		padding: 0;
		list-style: none;
	}

	li a {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: 16px;
		padding: 14px 22px;
		border-bottom: 1px solid var(--rule-soft);
		text-decoration: none;
	}

	li a:hover {
		background: var(--panel-sunk);
	}

	li a[aria-current='page'] {
		background: var(--ground);
	}

	li a[aria-current='page'] .name {
		font-weight: 600;
	}

	.name {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.short {
		color: var(--ink-2);
	}

	form {
		display: grid;
		gap: 8px;
		padding: 18px 22px 22px;
	}

	.row {
		display: flex;
		gap: 10px;
	}

	@media (max-width: 800px) {
		.split {
			grid-template-columns: minmax(0, 1fr);
			gap: 16px;
		}

		/* On a selected application, the detail owns the small screen; it links back. */
		.has-selection .list {
			display: none;
		}
	}
</style>
