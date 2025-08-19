<script lang="ts">
	import type { PageData } from './$types';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import CheckCircleIcon from '@lucide/svelte/icons/check-circle';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import PlayCircleIcon from '@lucide/svelte/icons/play-circle';
	import SkipForwardIcon from '@lucide/svelte/icons/skip-forward';
	import { invalidateAll } from '$app/navigation';
	import { cn } from '$lib/cn';

	let { data }: { data: PageData } = $props();

	let refreshing = $state(false);

	function formatDate(dateString: string | undefined): string {
		if (!dateString) return 'Never';
		const date = new Date(dateString);
		return new Intl.DateTimeFormat('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			hour12: false
		}).format(date);
	}

	function formatCommitMessage(message: string): string {
		return message.length > 60 ? message.substring(0, 60) + '...' : message;
	}

	function getStatusBadge(status: {
		is_active: boolean;
		is_skipped: boolean;
		error_message?: string;
		applied_at?: string;
	}) {
		if (status.error_message) {
			return { text: 'Error', class: 'bg-red-100 text-red-800', icon: AlertCircleIcon };
		}
		if (status.is_active) {
			return { text: 'Active', class: 'bg-blue-100 text-blue-800', icon: PlayCircleIcon };
		}
		if (status.is_skipped) {
			return { text: 'Skipped', class: 'bg-gray-100 text-gray-800', icon: SkipForwardIcon };
		}
		if (status.applied_at) {
			return { text: 'Completed', class: 'bg-green-100 text-green-800', icon: CheckCircleIcon };
		}
		return { text: 'Pending', class: 'bg-gray-100 text-gray-800', icon: ClockIcon };
	}

	async function handleRefresh() {
		refreshing = true;
		await invalidateAll();
		refreshing = false;
	}

	const activeMigration = $derived(data.statuses?.find((s) => s.is_active));
	const errorMigrations = $derived(data.statuses?.filter((s) => s.error_message) || []);
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-4xl font-medium">Sync Status</h1>
		<button
			onclick={handleRefresh}
			disabled={refreshing}
			class={cn(
				'flex items-center gap-2 rounded-md border border-gray-300 bg-white px-3 py-2 text-sm hover:bg-gray-50',
				refreshing && 'cursor-not-allowed opacity-50'
			)}
		>
			<RefreshCwIcon class={cn('h-4 w-4', refreshing && 'animate-spin')} />
			Refresh
		</button>
	</div>

	{#if data.error}
		<div class="rounded-md bg-red-50 p-4">
			<div class="flex">
				<AlertCircleIcon class="h-5 w-5 text-red-400" />
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800">Error loading sync data</h3>
					<div class="mt-2 text-sm text-red-700">
						<p>{data.error}</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	{#if activeMigration}
		<div class="rounded-md bg-blue-50 p-4">
			<div class="flex">
				<PlayCircleIcon class="h-5 w-5 text-blue-400" />
				<div class="ml-3">
					<h3 class="text-sm font-medium text-blue-800">Active Migration</h3>
					<div class="mt-2 text-sm text-blue-700">
						<p><strong>Commit:</strong> {activeMigration.commit.substring(0, 7)}</p>
						<p><strong>Message:</strong> {activeMigration.commit_message}</p>
						<p><strong>Date:</strong> {formatDate(activeMigration.commit_date)}</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	{#if errorMigrations.length > 0}
		{#each errorMigrations as errorMigration (errorMigration.commit)}
			<div class="rounded-md bg-red-50 p-4">
				<div class="flex">
					<AlertCircleIcon class="h-5 w-5 text-red-400" />
					<div class="ml-3 flex-1">
						<h3 class="text-sm font-medium text-red-800">Migration Error</h3>
						<div class="mt-2 text-sm text-red-700">
							<p>
								<strong>Commit:</strong>
								{errorMigration.commit.substring(0, 7)} - {errorMigration.commit_message}
							</p>
							<p class="mt-1"><strong>Error:</strong> {errorMigration.error_message}</p>
						</div>
					</div>
				</div>
			</div>
		{/each}
	{/if}

	<div class="space-y-6">
		<section>
			<h2 class="mb-3 text-xl font-medium">Recent Sync History</h2>
			<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
				<table class="w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Status
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Commit
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Message
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Commit Date
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Applied At
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200">
						{#if data.statuses && data.statuses.length > 0}
							{#each data.statuses as status (status.commit)}
								{@const badge = getStatusBadge(status)}
								<tr class="hover:bg-gray-50">
									<td class="px-6 py-3 whitespace-nowrap">
										<span
											class={cn(
												'inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-semibold',
												badge.class
											)}
										>
											<badge.icon class="h-3 w-3" />
											{badge.text}
										</span>
									</td>
									<td class="px-6 py-3 text-sm whitespace-nowrap text-gray-900">
										<a
											href={'https://github.com/' + data.repository + '/commit/' + status.commit}
											class="hover:underline"
											target="_blank"
										>
											<code class="rounded bg-gray-100 px-1 py-0.5 text-xs"
												>{status.commit.substring(0, 7)}</code
											>
										</a>
									</td>
									<td class="px-6 py-3 text-sm text-gray-500">
										{formatCommitMessage(status.commit_message)}
									</td>
									<td class="px-6 py-3 text-sm whitespace-nowrap text-gray-500">
										{formatDate(status.commit_date)}
									</td>
									<td class="px-6 py-3 text-sm whitespace-nowrap text-gray-500">
										{formatDate(status.applied_at)}
									</td>
								</tr>
							{/each}
						{:else}
							<tr>
								<td colspan="5" class="px-6 py-8 text-center text-sm text-gray-500">
									No sync history available
								</td>
							</tr>
						{/if}
					</tbody>
				</table>
			</div>
		</section>

		<section>
			<h2 class="mb-3 text-xl font-medium">Scheduled Jobs</h2>
			<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
				<table class="w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Job Name
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Schedule
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Status
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Last Run
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
							>
								Next Run
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200">
						{#if data.jobs && data.jobs.length > 0}
							{#each data.jobs as job (job.name)}
								<tr class="hover:bg-gray-50">
									<td class="px-6 py-3 text-sm font-medium text-gray-900">
										{job.name}
									</td>
									<td class="px-6 py-3 text-sm text-gray-500">
										<code class="rounded bg-gray-100 px-1 py-0.5 text-xs">{job.schedule}</code>
									</td>
									<td class="px-6 py-3 whitespace-nowrap">
										{#if job.is_running}
											<span
												class="inline-flex items-center gap-1 rounded-full bg-yellow-100 px-2 py-1 text-xs font-semibold text-yellow-800"
											>
												<RefreshCwIcon class="h-3 w-3 animate-spin" />
												Running
											</span>
										{:else}
											<span
												class="inline-flex items-center gap-1 rounded-full bg-gray-100 px-2 py-1 text-xs font-semibold text-gray-800"
											>
												<ClockIcon class="h-3 w-3" />
												Idle
											</span>
										{/if}
									</td>
									<td class="px-6 py-3 text-sm whitespace-nowrap text-gray-500">
										{formatDate(job.last_run)}
									</td>
									<td class="px-6 py-3 text-sm whitespace-nowrap text-gray-500">
										{formatDate(job.next_run)}
									</td>
								</tr>
							{/each}
						{:else}
							<tr>
								<td colspan="5" class="px-6 py-8 text-center text-sm text-gray-500">
									No scheduled jobs configured
								</td>
							</tr>
						{/if}
					</tbody>
				</table>
			</div>
		</section>
	</div>
</div>
