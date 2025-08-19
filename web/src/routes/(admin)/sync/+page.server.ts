import type { PageServerLoad } from './$types';
import type { SyncStatus, JobStatus } from '$lib/types/sync';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const [statusResponse, jobsResponse] = await Promise.all([
			fetch('/api/v1/sync/status?limit=10'),
			fetch('/api/v1/sync/jobs')
		]);

		if (!statusResponse.ok || !jobsResponse.ok) {
			return {
				statuses: [] as SyncStatus[],
				jobs: [] as JobStatus[],
				error: 'Failed to fetch sync data'
			};
		}

		const statuses: SyncStatus[] = await statusResponse.json();
		const jobs: JobStatus[] = await jobsResponse.json();

		return {
			statuses,
			jobs
		};
	} catch (error) {
		console.error('Error loading sync data:', error);
		return {
			statuses: [] as SyncStatus[],
			jobs: [] as JobStatus[],
			error: 'Failed to load sync status'
		};
	}
};