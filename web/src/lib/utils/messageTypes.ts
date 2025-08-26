export function getMessageType(form: any, message: string): 'success' | 'error' | 'warning' | 'info' {
	if (form.errors && Object.keys(form.errors).length > 0) {
        return 'error';
	}
	
	if (message.toLowerCase().includes('success') || 
		message.toLowerCase().includes('saved') || 
		message.toLowerCase().includes('created') ||
		message.toLowerCase().includes('updated')) {
		return 'success';
	}
	
	if (message.toLowerCase().includes('warning') ||
		message.toLowerCase().includes('exists')) {
		return 'warning';
	}
	
	if (message.toLowerCase().includes('error') ||
		message.toLowerCase().includes('failed') ||
		message.toLowerCase().includes('unauthorized') ||
		message.toLowerCase().includes('not found')) {
		return 'error';
	}
	
	return 'info';
}