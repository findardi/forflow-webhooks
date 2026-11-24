<script lang="ts">
	import { tick } from 'svelte';
	import { enhance } from '$app/forms';
	import { dateLong, dateShort } from '#lib/format.ts';
	import Messages from './Messages.svelte';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	let adding = $state(false);
	let copyLabel = $state('Copy ID');
	const addError = $derived(form && 'url' in form ? form.message : undefined);

	// Endpoint delete: two clicks. Confirm appears left of the trigger, Cancel takes its spot,
	// so a double click lands on Cancel.
	let armed = $state<string | null>(null);
	let deleting = $state<string | null>(null);

	async function disarm() {
		const id = armed;
		armed = null;
		await tick();
		document.getElementById(`del-${id}`)?.focus();
	}

	// Application delete: type the name. Cascades to every endpoint, so more friction.
	let confirmOpen = $state(false);
	let confirmName = $state('');
	let deletingApp = $state(false);
	const confirmError = $derived(form && 'confirm' in form ? form.message : undefined);
	const nameMatches = $derived(confirmName.trim() === data.application.name);
	const n = $derived(data.endpoints.length);

	async function closeConfirm() {
		confirmOpen = false;
		confirmName = '';
		await tick();
		document.getElementById('delete-app')?.focus();
	}

	const focus = (node: HTMLElement) => node.focus();

	let rotateFor = $state<string | null>(null);
	let rotating = $state(false);
	const rotateError = $derived(form && 'rotate' in form ? form.message : undefined);

	async function closeRotate() {
		const id = rotateFor;
		rotateFor = null;
		await tick();
		document.getElementById(`rot-${id}`)?.focus();
	}

	let dismissed = $state<string | null>(null);
	const reveal = $derived(form?.reveal && form.reveal.secret !== dismissed ? form.reveal : null);
	let secretCopyLabel = $state('Copy secret');

	async function copySecret(secret: string) {
		try {
			await navigator.clipboard.writeText(secret);
			secretCopyLabel = 'Copied';
		} catch {
			secretCopyLabel = 'Copy failed';
		}
		setTimeout(() => (secretCopyLabel = 'Copy secret'), 1500);
	}

	async function dismissReveal() {
		if (!reveal) return;
		const target = reveal.kind === 'created' ? 'endpoint-url' : `rot-${reveal.endpointId}`;
		dismissed = reveal.secret;
		await tick();
		document.getElementById(target)?.focus();
	}

	const selectAll = (node: HTMLInputElement) => {
		node.focus();
		node.select();
	};

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

<svelte:window
	onkeydown={(e) => {
		if (e.key !== 'Escape') return;
		if (armed) disarm();
		else if (rotateFor) closeRotate();
		else if (confirmOpen) closeConfirm();
	}}
/>

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
		<div class="actions">
			<button class="btn btn-ghost" onclick={copyId}>
				<span aria-live="polite">{copyLabel}</span>
			</button>
			<button
				id="delete-app"
				class="btn btn-ghost"
				aria-expanded={confirmOpen}
				aria-controls="delete-app-form"
				onclick={() => (confirmOpen ? closeConfirm() : (confirmOpen = true))}>Delete…</button
			>
		</div>
	</header>

	{#if confirmOpen}
		<form
			id="delete-app-form"
			class="danger"
			method="POST"
			action="?/deleteApplication"
			use:enhance={() => {
				deletingApp = true;
				return async ({ update }) => {
					await update();
					deletingApp = false;
				};
			}}
		>
			<h2 class="label">Delete application</h2>
			<p class="note">
				Deletes this application{n ? ` and its ${n} endpoint${n === 1 ? '' : 's'}` : ''}. This
				cannot be undone.
			</p>
			<label for="confirm-name">Type <strong>{data.application.name}</strong> to confirm</label>
			<div class="row">
				<input
					id="confirm-name"
					class="field"
					name="confirm"
					required
					autocomplete="off"
					spellcheck="false"
					bind:value={confirmName}
					aria-invalid={confirmError ? 'true' : undefined}
					aria-describedby={confirmError ? 'confirm-error' : undefined}
					{@attach focus}
				/>
				<button class="btn btn-danger" disabled={!nameMatches || deletingApp}>
					Delete application
				</button>
				<button type="button" class="btn btn-ghost" onclick={closeConfirm}>Cancel</button>
			</div>
			{#if confirmError}
				<p id="confirm-error" class="error">{confirmError}</p>
			{/if}
		</form>
	{/if}

	<nav class="tabs" aria-label="Application sections">
		<a
			href="/applications/{data.application.id}"
			aria-current={data.tab === 'endpoints' ? 'page' : undefined}
			>Endpoints <span class="count">{n}</span></a
		>
		<a href="?tab=messages" aria-current={data.tab === 'messages' ? 'page' : undefined}>Messages</a>
	</nav>

	{#if data.tab === 'messages'}
		<Messages
			appId={data.application.id}
			endpoints={data.endpoints}
			endpointId={data.endpointId}
			messages={data.messages}
		/>
	{:else}
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

		{#snippet secretReveal(r: { kind: 'created' | 'rotated'; url: string; secret: string })}
			<section class="reveal" aria-labelledby="reveal-title">
				<h2 id="reveal-title" class="label">
					{r.kind === 'created' ? 'Endpoint added' : 'Secret rotated'}
				</h2>
				{#if r.kind === 'created'}
					<p class="mono url">{r.url}</p>
				{/if}
				<label for="secret-value" class="secret-label">Signing secret</label>
				<div class="row">
					<input
						id="secret-value"
						class="field mono"
						readonly
						value={r.secret}
						autocomplete="off"
						spellcheck="false"
						onfocus={(e) => e.currentTarget.select()}
						aria-describedby="secret-note"
						{@attach selectAll}
					/>
					<button type="button" class="btn btn-ghost" onclick={() => copySecret(r.secret)}>
						<span aria-live="polite">{secretCopyLabel}</span>
					</button>
				</div>
				<p id="secret-note" class="note">
					{#if r.kind === 'created'}
						Shown once. Store it with your receiver to verify each webhook's signature. Forflow
						cannot show it again; rotate the secret to get a new one.
					{:else}
						Shown once. The previous secret keeps signing for 24 hours, so switch your receiver to
						this one before then.
					{/if}
				</p>
				<div class="buttons">
					<button type="button" class="text-btn" onclick={dismissReveal}>Done</button>
				</div>
			</section>
		{/snippet}

		{#if rotateError}
			<div class="strip">
				<p class="error">{rotateError}</p>
			</div>
		{/if}

		{#if reveal?.kind === 'created'}
			<div class="strip">{@render secretReveal(reveal)}</div>
		{/if}

		{#if data.endpoints.length}
			<div class="table-wrap">
				<table>
					<thead>
						<tr>
							<th class="label" scope="col">Endpoint URL</th>
							<th class="label id-col" scope="col">ID</th>
							<th class="label" scope="col">Created (UTC)</th>
							<th scope="col"><span class="sr-only">Actions</span></th>
						</tr>
					</thead>
					<tbody>
						{#each data.endpoints as ep (ep.id)}
							{@const rotatedHere = reveal?.kind === 'rotated' && reveal.endpointId === ep.id}
							<tr class:armed={armed === ep.id} class:open={rotateFor === ep.id || rotatedHere}>
								<td class="mono url">{ep.url}</td>
								<td class="mono dim id-col" title={ep.id}>{ep.id.slice(0, 8)}</td>
								<td class="time"
									><time datetime={ep.created_at}>{dateShort(ep.created_at)}</time></td
								>
								<td class="act">
									{#if armed === ep.id}
										<form
											method="POST"
											action="?/deleteEndpoint"
											class="confirm"
											use:enhance={() => {
												deleting = ep.id;
												return async ({ update }) => {
													await update();
													deleting = null;
													armed = null;
												};
											}}
										>
											<input type="hidden" name="endpoint_id" value={ep.id} />
											<button
												class="btn btn-s btn-danger"
												disabled={deleting === ep.id}
												aria-label="Confirm delete {ep.url}">Delete</button
											>
											<button
												type="button"
												class="btn btn-s btn-ghost"
												onclick={disarm}
												{@attach focus}>Cancel</button
											>
										</form>
									{:else}
										<button
											type="button"
											id="rot-{ep.id}"
											class="text-btn"
											aria-expanded={rotateFor === ep.id}
											aria-controls="rotate-{ep.id}"
											onclick={() => (rotateFor === ep.id ? closeRotate() : (rotateFor = ep.id))}
											aria-label="Rotate signing secret of {ep.url}">Rotate secret</button
										>
										<button
											type="button"
											id="del-{ep.id}"
											class="text-btn"
											onclick={() => (armed = ep.id)}
											aria-label="Delete endpoint {ep.url}">Delete</button
										>
									{/if}
								</td>
							</tr>
							{#if rotateFor === ep.id}
								<tr class="expand">
									<td colspan="4">
										<form
											id="rotate-{ep.id}"
											class="inline"
											method="POST"
											action="?/rotateSecret"
											use:enhance={() => {
												rotating = true;
												return async ({ update }) => {
													await update();
													rotating = false;
													rotateFor = null;
												};
											}}
										>
											<h2 class="label">Rotate signing secret</h2>
											<p class="note">
												A new secret starts signing now. The current one keeps signing beside it for
												24 hours, so move your receiver to the new secret before then. An older
												secret still inside its 24 hours stops signing immediately.
											</p>
											<input type="hidden" name="endpoint_id" value={ep.id} />
											<input type="hidden" name="url" value={ep.url} />
											<div class="buttons">
												<button class="btn btn-ghost" disabled={rotating} {@attach focus}
													>Rotate secret</button
												>
												<button type="button" class="text-btn" onclick={closeRotate}>Cancel</button>
											</div>
										</form>
									</td>
								</tr>
							{:else if rotatedHere && reveal}
								<tr class="expand">
									<td colspan="4">{@render secretReveal(reveal)}</td>
								</tr>
							{/if}
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

	/* Same voice as the top bar: current tab is graphite 600, the rest graphite 2. */
	.tabs {
		display: flex;
		gap: 28px;
		padding: 0 28px;
		border-bottom: 1px solid var(--rule);
	}

	.tabs a {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 48px;
		color: var(--ink-2);
		text-decoration: none;
	}

	.tabs a:hover {
		color: var(--ink);
	}

	.tabs a[aria-current='page'] {
		color: var(--ink);
		font-weight: 600;
	}

	.count {
		color: var(--ink-2);
		font-weight: 400;
		font-variant-numeric: tabular-nums;
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

	.act {
		width: 1%;
		padding-block: 9px;
		text-align: right;
		white-space: nowrap;
	}

	tbody tr.armed,
	tbody tr.armed:hover {
		background: var(--ground);
	}

	.text-btn {
		height: 32px;
		padding: 0 4px;
		border: 0;
		background: none;
		color: var(--ink-2);
		font: 500 14px var(--sans);
		cursor: pointer;
	}

	.text-btn:hover {
		color: var(--ink);
		text-decoration: underline;
		text-underline-offset: 3px;
	}

	.confirm {
		display: inline-flex;
		gap: 8px;
		padding: 0;
		border: 0;
	}

	.btn-s {
		height: 32px;
		padding: 0 12px;
		font-size: 14px;
	}

	.actions {
		display: flex;
		flex: none;
		gap: 10px;
	}

	/* Confirm panel sits right under the header, sunk so it reads as a temporary mode. */
	.danger {
		background: var(--panel-sunk);
	}

	.note {
		margin-bottom: 6px;
		color: var(--ink-2);
	}

	.strip {
		display: grid;
		gap: 8px;
		padding: 20px 28px;
		border-bottom: 1px solid var(--rule);
		background: var(--panel-sunk);
	}

	.reveal,
	.inline {
		display: grid;
		gap: 8px;
		padding: 0;
		border: 0;
	}

	.reveal .url {
		margin-bottom: 4px;
	}

	.reveal .note,
	.inline .note {
		max-width: 68ch;
		margin: 2px 0 0;
	}

	.expand .row {
		max-width: 760px;
	}

	tbody tr.open,
	tbody tr.open:hover,
	tbody tr.expand,
	tbody tr.expand:hover {
		background: var(--panel-sunk);
	}

	tbody tr.open td {
		border-bottom: 0;
	}

	.expand td {
		padding-top: 4px;
		padding-bottom: 24px;
	}

	.secret-label {
		color: var(--ink-2);
		font-weight: 500;
	}

	.buttons {
		display: flex;
		align-items: center;
		gap: 16px;
		margin-top: 6px;
	}

	.buttons > .text-btn:first-child {
		margin-left: -4px;
	}

	.act .text-btn + .text-btn {
		margin-left: 12px;
	}

	.danger strong {
		font-weight: 600;
		overflow-wrap: anywhere;
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

		form,
		.strip {
			padding: 16px;
		}

		.tabs {
			padding: 0 16px;
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

		.act .text-btn {
			display: block;
			margin-left: auto;
		}

		.act .text-btn + .text-btn {
			margin-left: auto;
		}
	}
</style>
