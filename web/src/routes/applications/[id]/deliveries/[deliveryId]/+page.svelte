<script lang="ts">
	import { dateExact, dateLong } from '#lib/format.ts';
	import type { Attempt } from '#lib/server/api.ts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const d = $derived(data.delivery);
	const n = $derived(data.attempts.length);

	const kinds: Record<string, string> = {
		dns: 'DNS lookup failed',
		timeout: 'Timed out',
		conn_refused: 'Connection refused',
		invalid_request: 'Invalid request',
		other: 'Network error'
	};

	const ok = (a: Attempt) =>
		a.response_status !== null && a.response_status >= 200 && a.response_status < 300;

	const result = (a: Attempt) =>
		a.response_status !== null
			? `HTTP ${a.response_status}`
			: (kinds[a.error_kind ?? 'other'] ?? a.error_kind);

	// What happens next depends on the status; a final delivery shows when it settled instead.
	const next = $derived.by(() => {
		if (d.status === 'processing') return { label: 'Next attempt', value: 'Sending now' };
		if (d.status === 'pending') {
			return d.attempt_count
				? { label: 'Next attempt (UTC)', value: dateExact(d.next_attempt_at) }
				: { label: 'Next attempt', value: 'Queued, not sent yet' };
		}
		return d.updated_at ? { label: 'Settled (UTC)', value: dateExact(d.updated_at) } : null;
	});

	const emptyText = $derived(
		d.status === 'canceled'
			? 'Canceled before the first attempt: its endpoint was deleted.'
			: d.status === 'processing'
				? 'The first attempt is in progress.'
				: 'No attempts yet. This delivery is queued for its first send.'
	);
</script>

<svelte:head>
	<title>Delivery {d.id.slice(0, 8)} · Forflow</title>
</svelte:head>

<section class="panel">
	<header>
		<div class="title">
			<h1>Delivery</h1>
			<p class="meta">
				<span class="mono">{d.id}</span>
				<span class="sep" aria-hidden="true">·</span>
				<span>Created {dateLong(d.created_at)}</span>
			</p>
		</div>
		<a class="btn btn-ghost back" href="/applications/{data.application.id}?tab=messages"
			>Back to messages</a
		>
	</header>

	<dl class="facts">
		<div>
			<dt class="label">Status</dt>
			<dd><span class="status {d.status}">{d.status}</span></dd>
		</div>
		<div class="wide">
			<dt class="label">Endpoint</dt>
			<dd class="mono url">{d.endpoint_url}</dd>
		</div>
		<div>
			<dt class="label">Attempts</dt>
			<dd class="num">{d.attempt_count}</dd>
		</div>
		{#if next}
			<div>
				<dt class="label">{next.label}</dt>
				<dd class="num">{next.value}</dd>
			</div>
		{/if}
	</dl>

	<div class="payload">
		<p class="label">Payload <span class="mono msg">message {d.message_id.slice(0, 8)}</span></p>
		<details>
			<summary class="mono">{JSON.stringify(d.payload)}</summary>
			<pre class="mono">{JSON.stringify(d.payload, null, 2)}</pre>
		</details>
	</div>

	<div class="list-head">
		<h2 class="label">Attempts</h2>
		<span class="label count">{n}</span>
	</div>

	{#if n}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th class="label num-col" scope="col">#</th>
						<th class="label" scope="col">Result</th>
						<th class="label time" scope="col">Time (UTC)</th>
					</tr>
				</thead>
				<tbody>
					{#each data.attempts as a, i (a.id)}
						<tr>
							<td class="num num-col">{n - i}</td>
							<td>
								<span class:fail={!ok(a)}>{result(a)}</span>
							</td>
							<td class="time"><time datetime={a.created_at}>{dateExact(a.created_at)}</time></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<div class="empty">
			<div class="grille" aria-hidden="true"></div>
			<p>{emptyText}</p>
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

	.title {
		min-width: 0;
	}

	h1 {
		font-size: 30px;
		font-weight: 600;
		line-height: 1.2;
	}

	.meta {
		display: flex;
		flex-wrap: wrap;
		gap: 4px 10px;
		margin-top: 6px;
		color: var(--ink-2);
	}

	.back {
		display: inline-flex;
		flex: none;
		align-items: center;
		text-decoration: none;
	}

	/* Facts read like the meta strip of an instrument: label over value, one row. */
	.facts {
		display: flex;
		flex-wrap: wrap;
		gap: 18px 40px;
		margin: 0;
		padding: 22px 28px;
		border-bottom: 1px solid var(--rule);
	}

	.facts div {
		display: grid;
		gap: 6px;
		min-width: 0;
	}

	.facts .wide {
		flex: 1 1 260px;
	}

	dd {
		margin: 0;
	}

	.url {
		overflow-wrap: anywhere;
	}

	.num {
		font-variant-numeric: tabular-nums;
	}

	.payload {
		display: grid;
		gap: 8px;
		padding: 20px 28px;
		border-bottom: 1px solid var(--rule);
	}

	/* Grid items default to min-width: auto; without this the summary pushes past the panel. */
	.payload > * {
		min-width: 0;
	}

	.msg {
		margin-left: 8px;
		color: var(--ink-3);
		font-weight: 400;
		letter-spacing: 0;
		text-transform: none;
	}

	summary {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		cursor: pointer;
	}

	summary:hover {
		text-decoration: underline;
		text-underline-offset: 3px;
	}

	details[open] summary {
		color: var(--ink-2);
	}

	pre {
		margin: 10px 0 0;
		padding: 12px 14px;
		overflow-x: auto;
		border: 1px solid var(--rule);
		border-radius: var(--radius-s);
		background: var(--panel-sunk);
		line-height: 1.5;
	}

	.list-head {
		display: flex;
		justify-content: space-between;
		padding: 18px 28px 0;
	}

	.count {
		font-variant-numeric: tabular-nums;
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

	.num-col {
		width: 1%;
		padding-right: 0;
		color: var(--ink-2);
	}

	.time {
		text-align: right;
		white-space: nowrap;
		font-variant-numeric: tabular-nums;
	}

	.empty {
		display: flex;
		align-items: center;
		gap: 20px;
		padding: 24px 28px 32px;
		color: var(--ink-2);
	}

	.empty .grille {
		flex: none;
		width: 72px;
		height: 40px;
	}

	@media (max-width: 800px) {
		header {
			flex-direction: column-reverse;
			gap: 16px;
			padding: 20px 16px;
		}

		h1 {
			font-size: 24px;
		}

		.meta {
			overflow-wrap: anywhere;
		}

		.sep {
			display: none;
		}

		.facts,
		.payload {
			padding: 16px;
		}

		.facts {
			gap: 14px 28px;
		}

		.list-head {
			padding: 16px 16px 0;
		}

		th,
		td {
			padding-inline: 16px;
		}
	}
</style>
