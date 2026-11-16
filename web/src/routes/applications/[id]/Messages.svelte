<script lang="ts">
	import { tick } from 'svelte';
	import { enhance } from '$app/forms';
	import { dateShort } from '#lib/format.ts';
	import type { DeliveredMessage, Endpoint } from '#lib/server/api.ts';

	type Props = {
		appId: string;
		endpoints: Endpoint[];
		endpointId: string | null;
		messages: DeliveredMessage[];
	};

	let { appId, endpoints, endpointId, messages }: Props = $props();

	const sample = '{\n  "type": "test.ping",\n  "data": { "hello": "world" }\n}';

	// The composer stays closed until asked for, so the tab opens on what was delivered.
	let open = $state(false);
	let draft = $state(sample);
	let sending = $state(false);
	let error = $state<string>();
	let sent = $state<string>();
	let trigger = $state<HTMLButtonElement>();

	async function close() {
		open = false;
		error = undefined;
		await tick();
		trigger?.focus();
	}

	const focus = (node: HTMLTextAreaElement) => node.focus();

	// Keep the tab and endpoint in the action URL; a bare ?/send would land back on Endpoints.
	const action = $derived(`?tab=messages${endpointId ? `&endpoint=${endpointId}` : ''}&/send`);
</script>

<div class="bar">
	{#if endpoints.length}
		<form method="GET" class="pick" data-sveltekit-keepfocus data-sveltekit-noscroll>
			<input type="hidden" name="tab" value="messages" />
			<label for="endpoint-pick" class="label">Endpoint</label>
			<select
				id="endpoint-pick"
				class="field mono"
				name="endpoint"
				value={endpointId ?? ''}
				onchange={(e) => e.currentTarget.form?.requestSubmit()}
			>
				<option value="">All endpoints</option>
				{#each endpoints as ep (ep.id)}
					<option value={ep.id}>{ep.url}</option>
				{/each}
			</select>
		</form>
	{/if}

	<div class="end">
		<p class="sent" role="status">
			{#if sent}Sent <span class="mono">{sent.slice(0, 8)}</span>{/if}
		</p>
		{#if !open}
			<button
				bind:this={trigger}
				class="btn btn-key"
				aria-controls="compose"
				aria-expanded="false"
				onclick={() => {
					open = true;
					sent = undefined;
				}}>Send message…</button
			>
		{/if}
	</div>
</div>

{#if open}
	<form
		id="compose"
		class="compose"
		method="POST"
		{action}
		use:enhance={() => {
			sending = true;
			return async ({ result, update }) => {
				await update({ reset: false });
				sending = false;
				if (result.type === 'failure') {
					error = String(result.data?.message ?? 'Could not send the message.');
				} else if (result.type === 'success') {
					sent = String(result.data?.sent ?? '');
					draft = sample;
					await close();
				}
			};
		}}
	>
		<label for="payload" class="label">Send message</label>
		<p class="note">
			Any JSON value. Forflow stores it and creates one delivery for each endpoint of this
			application.
		</p>
		<textarea
			id="payload"
			class="field mono"
			name="payload"
			rows="6"
			spellcheck="false"
			autocomplete="off"
			bind:value={draft}
			onkeydown={(e) => {
				if (e.key === 'Escape') close();
				if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) e.currentTarget.form?.requestSubmit();
			}}
			aria-invalid={error ? 'true' : undefined}
			aria-describedby={error ? 'payload-error' : 'payload-hint'}
			{@attach focus}></textarea>
		<div class="row">
			<button class="btn btn-key" disabled={sending}>Send message</button>
			<button type="button" class="btn btn-ghost" onclick={close}>Cancel</button>
			<span id="payload-hint" class="hint">Ctrl + Enter sends</span>
		</div>
		{#if error}
			<p id="payload-error" class="error">{error}</p>
		{/if}
	</form>
{/if}

{#snippet status(s: DeliveredMessage['delivery_status'])}
	{#if s}
		<span class="status {s}">{s}</span>
	{:else}
		<span class="status none" title="No endpoint existed when this message was sent">—</span>
	{/if}
{/snippet}

{#if !messages.length && !endpoints.length}
	<div class="empty">
		<div class="grille" aria-hidden="true"></div>
		<p>
			No endpoints yet, so a message has nowhere to go. It reaches only the endpoints that exist
			when it is sent. <a href="/applications/{appId}">Add an endpoint</a> first.
		</p>
	</div>
{:else if messages.length}
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th class="label" scope="col">Payload</th>
					{#if !endpointId}
						<th class="label endpoint-col" scope="col">Endpoint</th>
					{/if}
					<th class="label status-col" scope="col">Status</th>
					<th class="label id-col" scope="col">ID</th>
					<th class="label time" scope="col">Created (UTC)</th>
				</tr>
			</thead>
			<tbody>
				{#each messages as m (m.delivery_id ?? m.id)}
					<tr>
						<td class="payload">
							<details>
								<summary class="mono">{JSON.stringify(m.payload)}</summary>
								<pre class="mono">{JSON.stringify(m.payload, null, 2)}</pre>
							</details>
							<span class="inline">
								{@render status(m.delivery_status)}
								{#if m.endpoint_url}<span class="mono dim">{m.endpoint_url}</span>{/if}
							</span>
						</td>
						{#if !endpointId}
							<td class="mono endpoint-col" class:dim={!m.endpoint_url} title={m.endpoint_url}
								>{m.endpoint_url ?? '—'}</td
							>
						{/if}
						<td class="status-col">{@render status(m.delivery_status)}</td>
						<td class="mono dim id-col" title={m.id}>{m.id.slice(0, 8)}</td>
						<td class="time"><time datetime={m.created_at}>{dateShort(m.created_at)}</time></td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else}
	<div class="empty">
		<div class="grille" aria-hidden="true"></div>
		<p>
			{endpointId
				? 'No messages delivered to this endpoint yet.'
				: 'No messages for this application yet.'} Send one to see it here.
		</p>
	</div>
{/if}

<style>
	.bar {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 16px 24px;
		padding: 20px 28px;
		border-bottom: 1px solid var(--rule);
	}

	.pick {
		display: grid;
		flex: 0 1 460px;
		gap: 8px;
		min-width: 0;
	}

	select.field {
		width: 100%;
		cursor: pointer;
	}

	.end {
		display: flex;
		flex: none;
		align-items: center;
		gap: 16px;
		margin-left: auto;
	}

	.sent {
		color: var(--ink-2);
	}

	/* Sunk, like the delete strip: a temporary mode under the toolbar. */
	.compose {
		display: grid;
		gap: 8px;
		padding: 20px 28px;
		border-bottom: 1px solid var(--rule);
		background: var(--panel-sunk);
	}

	.note {
		margin-bottom: 6px;
		color: var(--ink-2);
	}

	textarea.field {
		height: auto;
		padding: 10px 12px;
		line-height: 1.5;
		resize: vertical;
		tab-size: 2;
	}

	.row {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.hint {
		margin-left: 6px;
		color: var(--ink-3);
		font-size: 14px;
	}

	.table-wrap {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		table-layout: fixed;
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
		vertical-align: top;
	}

	tbody tr:last-child td {
		border-bottom: 0;
	}

	tbody tr:hover {
		background: var(--panel-sunk);
	}

	.id-col {
		width: 140px;
	}

	.time {
		width: 170px;
		text-align: right;
		white-space: nowrap;
		font-variant-numeric: tabular-nums;
	}

	.dim {
		color: var(--ink-2);
	}

	/* Status is words, not colour: only a failure earns the error rust and its dot. */
	.status-col {
		width: 130px;
	}

	.endpoint-col {
		width: 240px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.status.none {
		color: var(--ink-3);
	}

	.status {
		text-transform: capitalize;
	}

	.status.pending {
		color: var(--ink-2);
	}

	.status.canceled {
		color: var(--ink-3);
	}

	.status.failed {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		color: var(--error);
		font-weight: 500;
	}

	.status.failed::before {
		content: '';
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--error-dot);
	}

	.inline {
		display: none;
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
		background: var(--panel);
		line-height: 1.5;
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
		.bar {
			flex-direction: column;
			align-items: stretch;
			padding: 16px;
		}

		.pick {
			flex-basis: auto;
		}

		.end {
			justify-content: space-between;
			margin-left: 0;
		}

		.compose {
			padding: 16px;
		}

		.hint,
		.id-col,
		.endpoint-col,
		.status-col {
			display: none;
		}

		.inline {
			display: flex;
			flex-wrap: wrap;
			gap: 2px 10px;
			margin-top: 4px;
			font-size: 14px;
			overflow-wrap: anywhere;
		}

		th,
		td {
			padding-inline: 16px;
		}

		.time {
			width: 130px;
		}
	}
</style>
