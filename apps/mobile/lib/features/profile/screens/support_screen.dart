import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class SupportScreen extends StatefulWidget {
  const SupportScreen({super.key});

  @override
  State<SupportScreen> createState() => _SupportScreenState();
}

class _SupportScreenState extends State<SupportScreen> {
  bool _isLoading = true;
  List<dynamic> _tickets = [];

  @override
  void initState() {
    super.initState();
    _loadTickets();
  }

  Future<void> _loadTickets() async {
    setState(() => _isLoading = true);
    try {
      final list = await ApiClient().get('/me/support-tickets');
      setState(() => _tickets = list ?? []);
    } catch (e) {
      debugPrint('Support error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showNewTicketDialog() {
    final subjectCtrl = TextEditingController();
    final descCtrl = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('문의 작성', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: subjectCtrl, decoration: const InputDecoration(labelText: '제목')),
            const SizedBox(height: 12),
            TextField(controller: descCtrl, maxLines: 3, decoration: const InputDecoration(labelText: '문의 내용')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(ctx);
              await ApiClient().post('/me/support-tickets', {
                'subject': subjectCtrl.text.trim(),
                'description': descCtrl.text.trim(),
                'category': 'GENERAL',
                'priority': 'NORMAL',
              });
              _loadTickets();
            },
            child: const Text('접수'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('1:1 고객지원'),
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _showNewTicketDialog),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Privacy note
                Container(
                  padding: const EdgeInsets.all(12),
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: AppColors.brandPrimary.withOpacity(0.08),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Text(
                    '계좌번호, 비밀번호 등 민감한 금융정보가 문의 내용에 포함되지 않도록 유의해 주세요.',
                    style: TextStyle(fontSize: 11, color: AppColors.brandPrimary),
                  ),
                ),

                if (_tickets.isEmpty)
                  const Center(
                    child: Padding(
                      padding: EdgeInsets.all(32.0),
                      child: Text('접수된 문의 내역이 없습니다.', style: TextStyle(color: AppColors.textSecondary)),
                    ),
                  )
                else
                  ..._tickets.map((t) => Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          title: Text(t['subject'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                          subtitle: Text(t['description'] ?? '', maxLines: 2, overflow: TextOverflow.ellipsis, style: const TextStyle(fontSize: 12)),
                          trailing: Text(
                            t['status'] ?? 'OPEN',
                            style: TextStyle(
                              color: t['status'] == 'RESOLVED' ? AppColors.success : AppColors.brandPrimary,
                              fontWeight: FontWeight.bold,
                              fontSize: 11,
                            ),
                          ),
                        ),
                      )),
              ],
            ),
    );
  }
}
