import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  Map<String, dynamic>? _user;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadProfile();
  }

  Future<void> _loadProfile() async {
    setState(() => _isLoading = true);
    try {
      final u = await ApiClient().get('/me');
      setState(() => _user = u);
    } catch (e) {
      debugPrint('Profile error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _requestDataExport() async {
    await ApiClient().post('/me/data-exports', {'format': 'CSV'});
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('CSV 데이터 내보내기가 생성되었습니다.')));
    }
  }

  Future<void> _requestAccountDeletion() async {
    if (await showDialog<bool>(
          context: context,
          builder: (ctx) => AlertDialog(
            title: const Text('회원 탈퇴 및 영구 삭제'),
            content: const Text('계정 삭제 시 모든 수동 금융 정보가 영구 파기됩니다. 계속하시겠습니까?'),
            actions: [
              TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('취소')),
              ElevatedButton(
                onPressed: () => Navigator.pop(ctx, true),
                style: ElevatedButton.styleFrom(backgroundColor: AppColors.danger),
                child: const Text('영구 삭제'),
              ),
            ],
          ),
        ) ==
        true) {
      await ApiClient().post('/me/deletion-requests', {});
      await ApiClient().clearToken();
      if (mounted) {
        context.go('/welcome');
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('설정 & 더보기')),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // User Card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Row(
                      children: [
                        CircleAvatar(
                          radius: 24,
                          backgroundColor: AppColors.brandPrimary.withOpacity(0.1),
                          child: Text(
                            (_user?['display_name'] ?? 'U').substring(0, 1),
                            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18, color: AppColors.brandPrimary),
                          ),
                        ),
                        const SizedBox(width: 16),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(_user?['display_name'] ?? '사용자', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                            Text(_user?['email'] ?? '', style: const TextStyle(color: AppColors.textSecondary, fontSize: 12)),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 20),

                // Menus
                _buildTile(Icons.headset_mic_outlined, '1:1 고객지원', () => context.push('/support')),
                _buildTile(Icons.download_outlined, '개인 데이터 내보내기 (CSV)', _requestDataExport),
                _buildTile(Icons.security, '보안 및 세션 확인', () {}),
                _buildTile(Icons.info_outline, '서비스 및 라이선스 정보', () {}),
                const Divider(height: 32),

                _buildTile(Icons.logout, '로그아웃', () async {
                  await ApiClient().clearToken();
                  if (context.mounted) context.go('/welcome');
                }),
                _buildTile(Icons.delete_forever, '회원 탈퇴 및 삭제 요청', _requestAccountDeletion, isDestructive: true),
              ],
            ),
    );
  }

  Widget _buildTile(IconData icon, String title, VoidCallback onTap, {bool isDestructive = false}) {
    return ListTile(
      leading: Icon(icon, color: isDestructive ? AppColors.danger : AppColors.textPrimary),
      title: Text(
        title,
        style: TextStyle(
          fontWeight: FontWeight.w600,
          fontSize: 14,
          color: isDestructive ? AppColors.danger : AppColors.textPrimary,
        ),
      ),
      trailing: const Icon(Icons.chevron_right, size: 18, color: AppColors.textSecondary),
      onTap: onTap,
    );
  }
}
