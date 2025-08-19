export interface SyncStatus {
	repo: string;
	commit: string;
	commit_message: string;
	commit_date: string;
	manifest?: string;
	applied_migration?: string;
	applied_at?: string;
	is_active: boolean;
	is_skipped: boolean;
	error_message?: string;
}

export interface JobStatus {
	name: string;
	schedule: string;
	last_run: string;
	next_run: string;
	is_running: boolean;
}
