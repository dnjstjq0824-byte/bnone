import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class NotificationsScreen extends StatefulWidget {
  const NotificationsScreen({super.key});

  @override
  State<NotificationsScreen> createState() => _NotificationsScreenState();
}

class _NotificationsScreenState extends State<NotificationsScreen> {
  bool _isLoading = true;
  List<dynamic> _notifications = [];

  @override
  void initState() {
    super.initState();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    setState(() => _isLoading = true);
    try {
      final list = await ApiClient().get('/me/notifications');
      setState(() => _notifications = list ?? []);
    } catch (e) {
      debugPrint('Notif error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _markRead(String id) async {
    await ApiClient().patch('/me/notifications/$id/read', {});
    _loadNotifications();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('알림함')),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _notifications.isEmpty
              ? const Center(child: Text('새로운 알림이 없습니다.', style: TextStyle(color: AppColors.textSecondary)))
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _notifications.length,
                  itemBuilder: (context, index) {
                    final n = _notifications[index];
                    final isRead = n['status'] == 'READ';

                    return Card(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: ListTile(
                        leading: CircleAvatar(
                          backgroundColor: isRead ? Colors.grey.withOpacity(0.1) : AppColors.brandPrimary.withOpacity(0.1),
                          child: Icon(
                            n['type'] == 'BUDGET_ALERT' ? Icons.warning_amber : Icons.notifications,
                            color: isRead ? Colors.grey : AppColors.brandPrimary,
                            size: 20,
                          ),
                        ),
                        title: Text(n['title'] ?? '', style: TextStyle(fontWeight: isRead ? FontWeight.normal : FontWeight.bold, fontSize: 14)),
                        subtitle: Text(n['body'] ?? '', style: const TextStyle(fontSize: 12)),
                        trailing: isRead
                            ? null
                            : IconButton(
                                icon: const Icon(Icons.check, size: 18),
                                onPressed: () => _markRead(n['id']),
                              ),
                      ),
                    );
                  },
                ),
    );
  }
}
