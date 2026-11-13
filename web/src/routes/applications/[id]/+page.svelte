<script lang="ts">
	import { enhance } from '$app/forms';
	import { dateLong, dateShort } from '#lib/format.ts';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	let adding = $state(false);
	let copyLabel = $state('Copy ID');
	const addError = $derived(form && 'url' in form ? form.message : undefined);

	async function copyId() {
		try {
			await navigator.clipboard.writeText(data.application.id);
			copyLabel = 'Copied';
		} catch {
			// navigator.clipboard is undefined on non-secure origins (plain http on a LAN host).
			copyLabel = 'Copy failed';
		}
		setTimeout(() => (copyLabel = 'Copy ID'), 1500);
	}
</script>

<svelte:head>
	<title>{data.application.name} · Forflow</title>
</svelte:head>

<section class="panel">
	<header>
		<a class="back" href="/applications">All applications</a>
		<div class="title">
			<h1>{data.application.name}</h1>
			<p class="meta">
				<span class="mono">{data.application.id}</span>
				<span class="sep" aria-hidden="true">·</span>
				<span>Created {dateLong(data.application.created_at)}</span>
			</p>
		</div>
		<button class="btn btn-ghost" onclick={copyId}>
			<span aria-live="polite">{copyLabel}</span>
		</button>
	</header>

	<form
		method="POST"
		action="?/add"
		use:enhance={() => {
			adding = true;
			return async ({ update }) => {
				await update();
				adding = false;
			};
		}}
	>
		<label for="endpoint-url" class="label">Add endpoint</label>
		<div class="row">
			<input
				id="endpoint-url"
				class="field mono"
				name="url"
				type="url"
				required
				autocomplete="off"
				spellcheck="false"
				placeholder="https://your-service.example.com/webhooks"
				value={addError ? (form?.url ?? '') : ''}
				aria-invalid={addError ? 'true' : undefined}
				aria-describedby={addError ? 'endpoint-url-error' : undefined}
			/>
			<button class="btn btn-key" disabled={adding}>Add endpoint</button>
		</div>
		{#if addError}
			<p id="endpoint-url-error" class="error">{addError}</p>
		{/if}
	</form>

	{#if data.endpoints.length}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th class="label" scope="col">Endpoint URL</th>
						<th class="label id-col" scope="col">ID</th>
						<th class="label" scope="col">Created (UTC)</th>
					</tr>
				</thead>
				<tbody>
					{#each data.endpoints as ep (ep.id)}
						<tr>
							<td class="mono url">{ep.url}</td>
							<td class="mono dim id-col" title={ep.id}>{ep.id.slice(0, 8)}</td>
							<td class="time"><time datetime={ep.created_at}>{dateShort(ep.created_at)}</time></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<div class="empty">
			<div class="grille" aria-hidden="true"></div>
			<p>No endpoints yet. Add the URL that should receive this application's webhooks.</p>
		</div>
	{/if}
</section>

<style>
	header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 24px;
		padding: 26px 28px;
		border-bottom: 1px solid var(--rule);
	}

	.back {
		display: none;
	}

	.title {
		min-width: 0;
	}

	h1 {
		font-size: 30px;
		font-weight: 600;
		line-height: 1.2;
		overflow-wrap: anywhere;
	}

	.meta {
		display: flex;
		flex-wrap: wrap;
		gap: 4px 10px;
		margin-top: 6px;
		color: var(--ink-2);
	}

	form {
		display: grid;
		gap: 8px;
		padding: 20px 28px;
		border-bottom: 1px solid var(--rule);
	}

	.row {
		display: flex;
		gap: 10px;
	}

	.table-wrap {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th {
		padding: 12px 28px;
		text-align: left;
		border-bottom: 1px solid var(--rule);
		white-space: nowrap;
	}

	td {
		padding: 15px 28px;
		border-bottom: 1px solid var(--rule-soft);
	}

	tbody tr:last-child td {
		border-bottom: 0;
	}

	tbody tr:hover {
		background: var(--panel-sunk);
	}

	th:last-child,
	.time {
		text-align: right;
	}

	.url {
		overflow-wrap: anywhere;
	}

	.dim {
		color: var(--ink-2);
	}

	.time {
		white-space: nowrap;
		font-variant-numeric: tabular-nums;
	}

	.empty {
		display: flex;
		align-items: center;
		gap: 20px;
		padding: 32px 28px;
		color: var(--ink-2);
	}

	.empty .grille {
		flex: none;
		width: 72px;
		height: 40px;
	}

	@media (max-width: 800px) {
		header {
			flex-wrap: wrap;
			padding: 20px 16px;
		}

		.back {
			display: block;
			flex-basis: 100%;
			color: var(--ink-2);
		}

		h1 {
			font-size: 24px;
		}

		form {
			padding: 16px;
		}

		.row {
			flex-direction: column;
		}

		th,
		td {
			padding-inline: 16px;
		}

		.sep,
		.id-col {
			display: none;
		}
	}
</style>
