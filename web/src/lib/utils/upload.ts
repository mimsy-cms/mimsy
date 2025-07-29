export type UploadProgress = {
	id: string;
	file: File;
	progress: number;
	status: 'uploading' | 'completed' | 'error';
	error?: string;
};

export type UploadOptions = {
	url: string;
	onProgress?: (uploadId: string, progress: number) => void;
	onStatusChange?: (uploadId: string, status: UploadProgress['status'], error?: string) => void;
};

/**
 * uploadFile uploads a file to the specified URL and tracks its progress.
 *
 * @param file The file to upload
 * @param uploadId The unique identifier for the upload
 * @param options The options for the upload
 * @returns A promise that resolves with the upload result
 */
export function uploadFile(
	formData: FormData,
	uploadId: string,
	options: UploadOptions
): Promise<void> {
	return new Promise((resolve) => {
		/**
		 * We use XMLHttpRequest to handle file uploads as we want to track the
		 * file upload progress.
		 */
		const xhr = new XMLHttpRequest();

		xhr.upload.addEventListener('progress', (e) => {
			if (e.lengthComputable) {
				const progress = (e.loaded / e.total) * 100;
				options.onProgress?.(uploadId, progress);
			}
		});

		xhr.addEventListener('load', () => {
			if (xhr.status === 200 || xhr.status === 201) {
				options.onStatusChange?.(uploadId, 'completed');
				resolve();
			} else {
				const error = `Upload failed: ${xhr.statusText}`;
				options.onStatusChange?.(uploadId, 'error', error);
				resolve();
			}
		});

		xhr.addEventListener('error', () => {
			options.onStatusChange?.(uploadId, 'error', 'An error occurred during the upload');
			resolve();
		});

		xhr.open('POST', options.url);
		xhr.send(formData);
	});
}

export function createUploadProgress(files: FileList | File[]): UploadProgress[] {
	return Array.from(files).map((file) => ({
		// We need a unique ID for each upload.
		id: crypto.randomUUID(),
		file,
		progress: 0,
		status: 'uploading' as const
	}));
}
