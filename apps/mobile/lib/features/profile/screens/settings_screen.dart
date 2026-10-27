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
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('CSV 데이터 내보내기가 생성되었습니다.')),
      );
    }
  }

  Future<void> _requestAccountDeletion() async {
    final confirmed = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '회원 탈퇴 및 데이터 영구 파기',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: AppColors.danger),
            ),
            const SizedBox(height: 10),
            const Text(
              '회원 탈퇴 시 등록하신 모든 금융 장부, 거래 내역, 예산 데이터가 영구적으로 파기되며 복구할 수 없습니다. 계속하시겠습니까?',
              style: TextStyle(fontSize: 14, color: AppColors.textSecondary, height: 1.5),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () => Navigator.pop(ctx, false),
                    child: const Text('취소'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(backgroundColor: AppColors.danger),
                    onPressed: () => Navigator.pop(ctx, true),
                    child: const Text('탈퇴 및 파기'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );

    if (confirmed == true) {
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
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('전체 서비스 & 설정')),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator(color: AppColors.brandPrimary))
          : ListView(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              children: [
                // User Profile Card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Row(
                      children: [
                        Container(
                          width: 48,
                          height: 48,
                          decoration: BoxDecoration(
                            color: AppColors.primaryLight,
                            shape: BoxShape.circle,
                          ),
                          child: Center(
                            child: Text(
                              (_user?['display_name'] ?? 'U').substring(0, 1),
                              style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 20, color: AppColors.brandPrimary),
                            ),
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                _user?['display_name'] ?? '사용자',
                                style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 17, color: AppColors.textPrimary),
                              ),
                              const SizedBox(height: 2),
                              Text(
                                _user?['email'] ?? '',
                                style: const TextStyle(color: AppColors.textTertiary, fontSize: 13),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 12),

                // Financial Services Card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Padding(
                          padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                          child: Text(
                            '금융 생활',
                            style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700, color: AppColors.textTertiary),
                          ),
                        ),
                        _buildSettingItem(Icons.account_balance_wallet_outlined, '수동 자산 & 계정 장부', () => context.push('/assets')),
                        _buildSettingItem(Icons.receipt_long_outlined, '소비 & 거래 내역', () => context.push('/transactions')),
                        _buildSettingItem(Icons.track_changes_outlined, '목표 예산 계획', () => context.push('/budgets')),
                        _buildSettingItem(Icons.repeat_rounded, '정기 구독 관리', () => context.push('/subscriptions')),
                        _buildSettingItem(Icons.flight_takeoff_rounded, '여행 경비 & 환율', () => context.push('/travel')),
                        _buildSettingItem(Icons.currency_exchange_rounded, '공인 참고 환율', () => context.push('/fx')),
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 12),

                // Privacy & Support Card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Padding(
                          padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                          child: Text(
                            '고객지원 및 개인정보',
                            style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700, color: AppColors.textTertiary),
                          ),
                        ),
                        _buildSettingItem(Icons.headset_mic_outlined, '고객센터 & 1:1 문의', () => context.push('/support')),
                        _buildSettingItem(Icons.download_rounded, '내 금융 데이터 CSV 내보내기', _requestDataExport),
                        _buildSettingItem(
                          Icons.logout_rounded,
                          '로그아웃',
                          () async {
                            await ApiClient().clearToken();
                            if (mounted) context.go('/welcome');
                          },
                        ),
                        _buildSettingItem(
                          Icons.person_remove_outlined,
                          '회원 탈퇴 및 영구 파기',
                          _requestAccountDeletion,
                          textColor: AppColors.danger,
                          iconColor: AppColors.danger,
                        ),
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 24),
              ],
            ),
    );
  }

  Widget _buildSettingItem(
    IconData icon,
    String title,
    VoidCallback onTap, {
    Color? textColor,
    Color? iconColor,
  }) {
    return ListTile(
      leading: Icon(icon, color: iconColor ?? AppColors.textSecondary, size: 22),
      title: Text(
        title,
        style: TextStyle(
          fontSize: 15,
          fontWeight: FontWeight.w600,
          color: textColor ?? AppColors.textPrimary,
        ),
      ),
      trailing: const Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
      onTap: onTap,
    );
  }
}
