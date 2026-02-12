<?php
// Enable login without password for SQLite
function adminer_object()
{
    class AdminerSoftware extends Adminer
    {
        function login($login, $password)
        {
            // Always allow login for SQLite (no password needed)
            return true;
        }

        function loginForm()
        {
            // Auto-fill the login form (but don't auto-submit)
            ?>
            <table cellspacing="0" class="layout">
                <tr>
                    <th>
                        <?php echo lang('System'); ?>
                    <td>
                        <input type="hidden" name="auth[driver]" value="sqlite">
                        SQLite 3
                <tr>
                    <th>
                        <?php echo lang('Database'); ?>
                    <td>
                        <input name="auth[db]" value="/data/local.db" autocapitalize="off">
            </table>
            <p><input type="submit" value="<?php echo lang('Login'); ?>">
                <?php
        }
    }

    return new AdminerSoftware;
}

include __DIR__ . '/adminer.php';
