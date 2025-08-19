import type { PageServerLoad } from './$types';
import type { SyncStatus, JobStatus } from '$lib/types/sync';

type SyncStatusResponse = {
	statuses: SyncStatus[];
	repository: string;
};

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const [statusResponse, jobsResponse, activeMigrationResponse] = await Promise.all([
			fetch('/api/v1/sync/status?limit=10'),
			fetch('/api/v1/sync/jobs'),
			fetch('/api/v1/sync/active-migration')
		]);

		if (!statusResponse.ok || !jobsResponse.ok || !activeMigrationResponse.ok) {
			return {
				statuses: [] as SyncStatus[],
				jobs: [] as JobStatus[],
				activeMigration: null,
				error: 'Failed to fetch sync data'
			};
		}

		const { statuses, repository }: SyncStatusResponse = await statusResponse.json();
		const jobs: JobStatus[] = await jobsResponse.json();
		const { active_migration }: { active_migration: SyncStatus | null } =
			await activeMigrationResponse.json();

		return {
			statuses,
			jobs,
			repository,
			activeMigration: active_migration
		};
	} catch (error) {
		console.error('Error loading sync data:', error);
		return {
			statuses: [] as SyncStatus[],
			jobs: [] as JobStatus[],
			activeMigration: null,
			error: 'Failed to load sync status'
		};
	}
};
